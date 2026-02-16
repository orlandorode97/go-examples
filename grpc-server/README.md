# gRPC Server on Google Kubernetes Engine (GKE)

Este documento describe los pasos completos para:

* Activar GKE
* Crear un cluster
* Crear un repositorio en Artifact Registry
* Construir y subir una imagen Docker
* Desplegar el servidor gRPC en Kubernetes
* Conectarse usando `grpcurl`

---

# 1. Prerrequisitos

Instalar:

* Google Cloud SDK (`gcloud`)
* Docker
* kubectl
* grpcurl

Verificar:

```bash
gcloud version
docker version
kubectl version --client
grpcurl -version
```

Login:

```bash
gcloud auth login
```

Seleccionar proyecto:

```bash
gcloud config set project [PROJECT_ID]
```

---

# 2. Activar APIs necesarias

```bash
gcloud services enable \
  container.googleapis.com \
  artifactregistry.googleapis.com \
  compute.googleapis.com
```

---

# 3. Crear cluster GKE

```bash
gcloud container clusters create grpc-cluster \
  --zone us-central1-a \
  --num-nodes 1
```

Configurar acceso:

```bash
gcloud container clusters get-credentials grpc-cluster \
  --zone us-central1-a
```

Verificar:

```bash
kubectl get nodes
```

---

# 4. Crear repositorio Artifact Registry

```bash
gcloud artifacts repositories create my-repo \
  --repository-format=docker \
  --location=us-central1
```

Configurar autenticación Docker:

```bash
gcloud auth configure-docker us-central1-docker.pkg.dev
```

---

# 5. Construir y subir imagen Docker

Desde el directorio del proyecto:

```bash
docker buildx build \
  --platform linux/amd64 \
  -t us-central1-docker.pkg.dev/[PROJECT_ID]/my-repo/grpc-server:v1 \
  --push .
```

---

# 6. Crear deployment.yaml

Crear archivo `deployment.yaml`:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: grpc-server
spec:
  replicas: 1
  selector:
    matchLabels:
      app: grpc-server
  template:
    metadata:
      labels:
        app: grpc-server
    spec:
      containers:
        - name: grpc-server
          image: us-central1-docker.pkg.dev/[PROJECT_ID]/my-repo/grpc-server:v1
          ports:
            - containerPort: 50051

---
apiVersion: v1
kind: Service
metadata:
  name: grpc-server
spec:
  type: LoadBalancer
  selector:
    app: grpc-server
  ports:
    - name: grpc
      port: 50051
      targetPort: 50051
```

---

# 7. Desplegar en Kubernetes

```bash
kubectl apply -f deployment.yaml
```

Verificar pods:

```bash
kubectl get pods
```

Verificar service:

```bash
kubectl get svc grpc-server
```

Esperar hasta que aparezca EXTERNAL-IP:

```bash
kubectl get svc grpc-server
```

Ejemplo resultado:

```
NAME          TYPE           CLUSTER-IP     EXTERNAL-IP     PORT(S)
grpc-server   LoadBalancer   10.52.1.10     34.135.71.218   50051/TCP
```

---

# 8. Verificar endpoints

```bash
kubectl get endpoints grpc-server
```

---

# 9. Conectarse usando grpcurl

Listar servicios:

```bash
grpcurl -plaintext [EXTERNAL-IP]:50051 list
```

Ejemplo:

```bash
grpcurl -plaintext 34.135.71.218:50051 list
```

Output esperado:

```
generated.HelloService
```

---

# 10. Llamar método gRPC

```bash
grpcurl -plaintext \
  -d '{"name":"Orlando"}' \
  [EXTERNAL-IP]:50051 \
  generated.HelloService/SayHello
```

Respuesta esperada:

```json
{
  "message": "Hello, Orlando"
}
```

---

# 11. Ver logs

```bash
kubectl logs -l app=grpc-server
```

---

# 12. Arquitectura final

```
grpcurl client
   ↓
GCP LoadBalancer
   ↓
Kubernetes Service
   ↓
Pod
   ↓
Container
   ↓
Go gRPC server
```

---

# 13. Comandos útiles

Ver recursos:

```bash
kubectl get all
```

Reiniciar deployment:

```bash
kubectl rollout restart deployment grpc-server
```

Eliminar deployment:

```bash
kubectl delete -f deployment.yaml
```

Eliminar cluster:

```bash
gcloud container clusters delete grpc-cluster --zone us-central1-a
```

---

# Fin
