#!/bin/sh

params_num=$#
plugin_name=$1
image=$2
platform=$3

validate_param(){
if [ $params_num != 3 ];then
 cat <<EOF
err:  parameter transmission error, please refer to the following instructions:
usage: sh gen_plugin_yaml_templete.sh chckpr checkpr:latest gitee
params:
 - plugin name e.g: checkpr.
 - image  the image full name for depolyment container e.g: checkpr:latest.
 - platform the plugin for code platform e.g: gitee.
EOF
exit 1
fi
}

# make dir
create_plugin_dir(){
 if [ ! -d "$plugin_name" ];then
     mkdir "$plugin_name"
 else
    rm -rf "$plugin_name"
    mkdir "$plugin_name"
 fi
}

# create depolyment
create_depolyment(){
d_name="$plugin_name/depolyment.yaml"
cat > "${d_name}" <<EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: $plugin_name
  name: $plugin_name
spec:
  replicas: 1
  selector:
    matchLabels:
      component: $plugin_name
  template:
    metadata:
      labels:
        component: $plugin_name
    spec:
      containers:
        - args:
            - --$platform-token-path=/app/conf/robot-$platform/robot_token
            - --plugin-config=/app/conf/robot-$platform/$plugin_name.yaml
          image: $image
          imagePullPolicy: IfNotPresent
          name: $plugin_name-pod
          ports:
            - containerPort: 80
              name: http
              protocol: TCP
          volumeMounts:
            - mountPath: /app/conf/robot-$platform/robot_token
              subPath: $platform-token
              name: secret-volume
            - mountPath: app/conf/robot-$platform/$plugin_name.yaml
              subPath: $plugin_name.yaml
              name: $platform-$plugin_name
              readOnly: true
      volumes:
        - name: secret-volume
          secret:
            secretName: robot-server-secrets
        - name: $platform-$plugin_name
          configMap:
            name: $platform-$plugin_name
EOF
}

# create service yaml
create_service(){
svc_name="$plugin_name/service.yaml"
cat > "${svc_name}" <<EOF
apiVersion: v1
kind: Service
metadata:
  name: $plugin_name-service
spec:
  ports:
    - name: httpport
      port: 80
      protocol: TCP
      targetPort: 80
  selector:
    component: $plugin_name
  type: ClusterIP
EOF
}

# create configmap yaml file
create_configmap(){
cmap_name="$plugin_name/configmap.yaml"
cat > "${cmap_name}" <<EOF
apiVersion: v1
kind: ConfigMap
metadata:
  name: $platform-$plugin_name
data:
  $plugin_name.yaml: |
   #fill in the plugin config data
EOF
}

# cretate secrets yaml file
create_sec(){
sec_name="$plugin_name/secrets.yaml"
cat > "${sec_name}" <<EOF
apiVersion: secrets-manager.tuenti.io/v1alpha1
kind: SecretDefinition
metadata:
  name: robot-server-secrets
  spec:
    name: robot-server-secrets
    keysMap:
      $plugin_name-token:
    path: secrets/data/infra-common/robot-server
    key: $plugin_name-token
EOF
}

validate_param

create_plugin_dir

create_depolyment

create_service

create_configmap

create_sec
