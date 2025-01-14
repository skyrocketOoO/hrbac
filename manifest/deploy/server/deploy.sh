#!/bin/bash

# Variables
VM_USER="ubuntu"
VM_HOST="ec2-54-255-187-225.ap-southeast-1.compute.amazonaws.com"
PRIVATE_KEY="~/.ssh/vortex-dev.pem"
PRIVATE_KEY_PATH=$(eval echo $PRIVATE_KEY)
BINARY="server"

ssh -i "$PRIVATE_KEY" -T $VM_USER@$VM_HOST << ENDSSH
  sudo service alarmsystem stop
ENDSSH

make build
if [ $? -ne 0 ]; then
  echo "Error: Build failed."
  exit 1
fi

scp -i "$PRIVATE_KEY" "$BINARY" $VM_USER@$VM_HOST:~/app/
if [ $? -ne 0 ]; then
  echo "Error: File transfer failed."
  exit 1
fi

scp -i "$PRIVATE_KEY" ".env" $VM_USER@$VM_HOST:~/app/
if [ $? -ne 0 ]; then
  echo "Error: File transfer failed."
  exit 1
fi

# Connect to the VM and set up the project
ssh -i "$PRIVATE_KEY" -T $VM_USER@$VM_HOST << ENDSSH
  sudo service alarmsystem restart
ENDSSH

rm "$BINARY"
if [ $? -ne 0 ]; then
  echo "Error: Failed to clean up local binary."
  exit 1
fi

tag=$(date +'%Y-%m-%d-%H-%M-%S')
git tag "$tag"
git push origin "$tag"

echo "deployment complete!"