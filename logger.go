package main

import (
	"context"
	"encoding/json"
	"math"
	"os"
	"time"

	"github.com/gregtuc/docker-data-metrics/database"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

var (
	cli            *client.Client
	err            error
	previousCPU    uint64
	previousSystem uint64
)

func StartLogging() {

	//Create a new docker client
	ctx := context.Background()
	cli, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	//Initialize a list of Docker containers
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	//Get the container id by name
	containerId := getContainerId("/"+os.Getenv("CONTAINER_NAME"), containers)

	ticker := time.NewTicker(1 * time.Second)
	for range ticker.C {

		//Get container stats once
		rawStats, err := cli.ContainerStatsOneShot(ctx, containerId)
		if err != nil {
			panic(err)
		}

		//Decode the container stats
		var decodedStats *types.StatsJSON
		json.NewDecoder(rawStats.Body).Decode(&decodedStats)

		//Get container ip address
		containerConfig, err := cli.ContainerInspect(ctx, containerId)
		if err != nil {
			panic(err)
		}
		containerIp := containerConfig.NetworkSettings.IPAddress

		//Convert timestamp to unix
		timestamp := decodedStats.Read.Unix()

		//Fetch CPU & RAM usage; convert to percent
		var cpuPercent float64
		var ramPercent float64
		if rawStats.OSType != "windows" {
			cpuPercent = calculateCPUPercentUnix(previousCPU, previousSystem, decodedStats)
			previousCPU = decodedStats.CPUStats.CPUUsage.TotalUsage
			previousSystem = decodedStats.CPUStats.SystemUsage
		} else {
			cpuPercent = calculateCPUPercentWindows(decodedStats)
		}
		ramPercent = float64(decodedStats.MemoryStats.Usage) / float64(decodedStats.MemoryStats.Limit) * 100

		//Write to database
		go database.Write(database.Log{
			ContainerIP: containerIp,
			Timestamp:   timestamp,
			CPUPercent:  math.Round(cpuPercent*100) / 100,
			RamPercent:  math.Round(ramPercent*100) / 100,
		})

	}
}

//Get container id for a given container name
func getContainerId(containerName string, containers []types.Container) string {
	for _, container := range containers {
		if container.Names[0] == containerName {
			containerName = container.ID
			return container.ID
		}
	}
	return ""
}

/*
	Source: Helper conversion function that was used in the official Docker CLI
	https://github.com/docker/cli/blob/a32cd16160f1b41c1c4ae7bee4dac929d1484e59/cli/command/container/stats_helpers.go
*/
func calculateCPUPercentUnix(previousCPU, previousSystem uint64, v *types.StatsJSON) float64 {
	var (
		cpuPercent = 0.0
		// calculate the change for the cpu usage of the container in between readings
		cpuDelta = float64(v.CPUStats.CPUUsage.TotalUsage) - float64(previousCPU)
		// calculate the change for the entire system between readings
		systemDelta = float64(v.CPUStats.SystemUsage) - float64(previousSystem)
	)

	if systemDelta > 0.0 && cpuDelta > 0.0 {
		cpuPercent = (cpuDelta / systemDelta) * 100.0
	}
	return cpuPercent
}

/*
	Source: Helper conversion function that was used in the official Docker CLI
	https://github.com/docker/cli/blob/a32cd16160f1b41c1c4ae7bee4dac929d1484e59/cli/command/container/stats_helpers.go
*/
func calculateCPUPercentWindows(v *types.StatsJSON) float64 {
	// Max number of 100ns intervals between the previous time read and now
	possIntervals := uint64(v.Read.Sub(v.PreRead).Nanoseconds()) // Start with number of ns intervals
	possIntervals /= 100                                         // Convert to number of 100ns intervals
	possIntervals *= uint64(v.NumProcs)                          // Multiple by the number of processors
	// Intervals used
	intervalsUsed := v.CPUStats.CPUUsage.TotalUsage - v.PreCPUStats.CPUUsage.TotalUsage
	// Percentage avoiding divide-by-zero
	if possIntervals > 0 {
		return float64(intervalsUsed) / float64(possIntervals) * 100.0
	}
	return 0.00
}
