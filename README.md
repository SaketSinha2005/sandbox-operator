# Kubernetes Sandbox Execution Operator

A Kubernetes Operator built in Go using Operator SDK that manages isolated code execution environments as first-class cluster resources. Define a `SandboxEnvironment` custom resource and the operator handles the full lifecycle — provisioning, execution, status tracking, timeout enforcement, and cleanup.

---

## Overview

The operator introduces a `SandboxEnvironment` CRD to the cluster. Each CR describes a sandboxed execution unit: the runtime image, resource limits, storage, network access, and security constraints. The controller reconciles the desired state by managing a Pod and an optional PersistentVolumeClaim, tracking status transitions from provisioning through to completion or timeout.

---

## Architecture

```
SandboxEnvironment CR
        |
        v
  SandboxEnvironmentReconciler
        |
        |-- reconcilePVC()   --> PersistentVolumeClaim (workspace storage)
        |-- reconcilePod()   --> Pod (sandboxed container)
        |-- syncStatus()     --> CR status subresource
        |-- handleDeletion() --> finalizer cleanup
```

Owned resources (Pod and PVC) carry owner references back to the CR, so deleting a `SandboxEnvironment` garbage-collects all child resources automatically.

---

## CRD Spec Reference

```yaml
apiVersion: sandbox.mylab.io/v1alpha1
kind: SandboxEnvironment
metadata:
  name: my-sandbox
spec:
  runtime:
    image: python:3.12          # Container image to run
    language: python            # Enum: python | cpp | java
    command: [python3, -c, "print('hello')"]

  resources:
    requests:
      cpu: "250m"
      memory: "256Mi"
    limits:
      cpu: "1"
      memory: "512Mi"

  storage:
    size: "1Gi"                 # PVC size (format: <n>Gi)
    mountPath: "/workspace"     # Mount path inside the container

  network:
    enabled: false              # Labels pod for NetworkPolicy enforcement

  security:
    runAsNonRoot: true          # Enforced via SecurityContext
    readOnlyRootFilesystem: true
    allowPrivilegeEscalation: false

  timeout: 120s                 # activeDeadlineSeconds + operator-level check
```

---

## Status Fields

| Field     | Description                                              |
|-----------|----------------------------------------------------------|
| `phase`   | `Provisioning`, `Running`, `Succeeded`, `Failed`, `TimedOut` |
| `ready`   | `true` when the sandbox pod is actively running          |
| `podName` | Name of the managed Pod                                  |
| `podIP`   | IP assigned to the Pod                                   |
| `message` | Human-readable status detail                             |

---

## Security Model

- Pods run as UID 1000 (`runAsNonRoot: true` + `runAsUser: 1000`)
- Root filesystem is read-only; `/tmp` is provided via `emptyDir` for runtime writability
- Privilege escalation is disabled at the container level
- Network isolation is signaled via the `sandbox.mylab.io/network-disabled: "true"` label, intended for use with a companion `NetworkPolicy`
- `RestartPolicy: Never` ensures sandboxes run exactly once

---

## Timeout Enforcement

Timeout is enforced at two layers:

1. `activeDeadlineSeconds` on the Pod spec — Kubernetes kills the pod at the deadline
2. Operator-level check — if the CR has been in `Running` phase past the timeout duration, the operator transitions the phase to `TimedOut` and requeues accordingly

---

## Getting Started

**Prerequisites**

- Go 1.21+
- Operator SDK v1.34.1
- A running cluster (minikube, kind, or remote)
- `kubectl` configured

**Install CRD and run locally**

```bash
make generate
make manifests
kubectl apply -f config/crd/bases/

make run
```

**Apply a sample sandbox**

```bash
kubectl apply -f config/samples/sandbox_v1alpha1_sandboxenvironment.yaml
```

**Watch status**

```bash
kubectl get sandboxenvironment sandboxenvironment-sample -o jsonpath='{.status}' | jq
kubectl get pod,pvc | grep sandbox
```

**Delete and clean up**

```bash
kubectl delete sandboxenvironment sandboxenvironment-sample
# Pod and PVC are garbage-collected automatically
```

---

## Project Structure

```
sandbox-operator/
├── api/v1alpha1/
│   ├── sandboxenvironment_types.go     # CRD type definitions
│   └── groupversion_info.go
├── internal/controller/
│   └── sandboxenvironment_controller.go  # Reconciler implementation
├── config/
│   ├── crd/bases/                      # Generated CRD manifests
│   └── samples/                        # Example CR
└── cmd/main.go                         # Operator entrypoint
```

---

## RBAC

The operator requires the following permissions, declared via kubebuilder markers:

- `sandboxenvironments`: get, list, watch, create, update, patch, delete
- `sandboxenvironments/status`: get, update, patch
- `sandboxenvironments/finalizers`: update
- `pods`: get, list, watch, create, update, patch, delete
- `pods/log`: get, list
- `persistentvolumeclaims`: get, list, watch, create, update, patch, delete
