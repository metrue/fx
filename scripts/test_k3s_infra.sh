#!/bin/bash
set -e

master_ip=$(multipass list | tail -3 | grep k3s-master | awk -F' '  '{print $3}')
agent_1_ip=$(multipass list | tail -3 | grep k3s-worker1 | awk -F' '  '{print $3}')
agent_2_ip=$(multipass list | tail -3 | grep k3s-worker2 | awk -F' '  '{print $3}')
user="multipass"
echo SSH_KEY_FILE=./test/id_rsa ./build/fx infra create -name k3s-test-cloud -t k3s --master ${user}@${master_ip}  #--agents ${user}@${agent_1_ip},${user}@${agent_2_ip}
SSH_KEY_FILE=./test/id_rsa ./build/fx infra create -name k3s-test-cloud -t k3s --master ${user}@${master_ip}  #--agents ${user}@${agent_1_ip},${user}@${agent_2_ip}
