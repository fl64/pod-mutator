# pod-mutator webhook

```shell
operator-sdk init --plugins go/v3  --owner "fl64 <flsixtyfour@gmail.com>"
operator-sdk create api --group=core --version=v1 --kind=Pod --controller=true --resource=false 
make
operator-sdk create webhook --group core --version v1 --kind Pod --defaulting
make generate
make docker-build docker-push IMG=fl64/pod-mutator:v0.0.1
make deploy IMG=fl64/pod-mutator:v0.0.1
make undeploy 
```