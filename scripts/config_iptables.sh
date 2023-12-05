#!/bin/bash

# Set forwarding
sysctl -w net.ipv4.ip_forward=1

# Set rules for network interface
iptables -t nat -A POSTROUTING -o <network interface name> -j MASQUERADE
iptables -A FORWARD -p tcp -m tcp --tcp-flags SYN,RST SYN -j TCPMSS --set-mss 1400

# Clear existing rules
iptables -F
iptables -P INPUT DROP
iptables -P FORWARD DROP
iptables -P OUTPUT DROP

# Allow loopback interface
iptables -A INPUT -i lo -j ACCEPT
iptables -A OUTPUT -o lo -j ACCEPT

# Allow established and related connections
iptables -A INPUT -m conntrack --ctstate ESTABLISHED,RELATED -j ACCEPT
iptables -A OUTPUT -m conntrack --ctstate ESTABLISHED -j ACCEPT

# Allow incoming ICMP (ping) traffic
iptables -A INPUT -p icmp --icmp-type echo-request -j ACCEPT
iptables -A OUTPUT -p icmp --icmp-type echo-reply -j ACCEPT

# Allow outgoing ICMP (ping) traffic
iptables -A OUTPUT -p icmp --icmp-type echo-request -j ACCEPT
iptables -A INPUT -p icmp --icmp-type echo-reply -j ACCEPT

# Allow incoming DNS responses (UDP on port 53)
iptables -A INPUT -p udp --sport 53 -j ACCEPT
iptables -A OUTPUT -p udp --dport 53 -j ACCEPT

# Allow outgoing DNS requests (UDP on port 53)
iptables -A OUTPUT -p udp --dport 53 -j ACCEPT

# Allow outgoing HTTP and HTTPS traffic (TCP on ports 80 and 443)
iptables -A OUTPUT -p tcp --dport 80 -j ACCEPT
iptables -A OUTPUT -p tcp --dport 443 -j ACCEPT

# Allow outgoing traffic to external websites
iptables -A OUTPUT -o enps03 -j ACCEPT

# Save iptables rules
iptables-save > /etc/iptables/rules.v4

# Display current rules
iptables -L -n

# Stop firewall
systemctl stop ufw

echo "Iptables rules configured and saved successfully."
