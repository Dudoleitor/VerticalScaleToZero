# VerticalScaleToZero

## Description

**VerticalScaleToZero** is a tool developed as part of a research project with **Politecnico di Milano**. The tool is designed to show a proof of concept of scaling the resources of Kubernetes pods to zero, enabling resource optimization and cost reduction, especially in cloud environments.

The aim of the project is to investigate scalable, cost-effective strategies for managing Kubernetes-based workloads in cloud environments, with a focus on minimizing resource consumption. This project is built around a KIND (Kubernetes IN Docker) cluster, a tool that runs Kubernetes clusters in Docker containers. KIND is an ideal choice for local development and testing environments, as it allows the creation of lightweight, easily managed Kubernetes clusters.

To use VerticalScaleToZero, the **InPlacePodVerticalScaling** feature gate must be enabled in your Kubernetes cluster. This feature gate is essential for managing vertical pod scaling on existing pods (without the need for pod restarts).

## Prototype Status

It is important to note that VerticalScaleToZero is still in its prototype phase and was developed as part of a research project. The tool demonstrates the concept of workload scaling to zero replicas within Kubernetes, but it is not yet intended for widespread production use. It may lack certain features and robustness that would be expected from a final product.

The tool’s main goal is to contribute insights to the research on cloud resource optimization. Users should consider it a proof-of-concept for experimental or academic purposes.
For more detailed information about the research and the methodologies used, one can read the full **research report**:

- [Research Report: Vertical Scale to Zero for Kubernetes Pods](https://github.com/Dudoleitor/VerticalScaleToZero/blob/main/VerticalScaleToZero-report.pdf)

## Requirements

This tool is designed to work with **KIND** (Kubernetes IN Docker), so you will need to have KIND installed on your machine. KIND is a tool for running Kubernetes clusters locally using Docker containers, and it is an ideal choice for quickly setting up isolated Kubernetes clusters for testing and experimentation.

### Install KIND

Follow the official installation instructions for KIND:

- [KIND Installation Guide](https://kind.sigs.k8s.io/docs/user/quick-start/)


## Scripts for Environment Setup and Teardown

Two helper scripts are provided to assist in preparing and shutting down the environment.
#### `build_everything.sh`

This script automatically creates a KIND cluster with the InPlacePodVerticalScaling feature gate enabled and sets up all necessary resources for the tool.

Usage:

```bash
./build_everything.sh
```

### `destroy_everything.sh`

This script will tear down the environment by deleting the KIND cluster, ensuring a clean state.

Usage:

```bash
./destroy_everything.sh
```
