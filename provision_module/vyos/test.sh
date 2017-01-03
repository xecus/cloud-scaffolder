#!/bin/bash

export LC_ALL=C.UTF-8
python main.py provision \
	--target_router_ip 192.168.33.9 \
	--target_router_port 22 \
	--target_router_username vagrant \
	--target_router_password vagrant \
	--router_hostname vyos-taguro \
	--dhcp_ip_start 192.168.34.100 \
	--dhcp_ip_stop 192.168.34.150 \
	--dhcp_subnet 192.168.34.0/24 \
	--dhcp_default_dns 8.8.8.8 \
	--dhcp_default_router 192.168.34.1
