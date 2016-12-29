# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|

  config.vm.define "router" do |d|
    d.vm.hostname = "router"
    d.vm.box = "higebu/vyos"
    d.vm.network "private_network", ip: "192.168.33.9"
    d.vm.network "private_network", ip: "192.168.34.1", virtualbox__intnet: "hogenet"
    d.vm.network :forwarded_port, host: 3022, guest: 22
    config.vm.provider :virtualbox do |vb|
      vb.customize ["modifyvm", :id, "--memory", "2048", "--cpus", "2", "--ioapic", "on"]
    end
    config.vm.provision "shell", inline: <<-SHELL
    SHELL
  end

  config.vm.define "main" do |d|
    d.vm.hostname = "main"
    d.vm.box = "ubuntu/trusty64"
    d.vm.network "private_network", ip: "192.168.33.10"
    d.vm.network "private_network", ip: "192.168.34.0", virtualbox__intnet: "hogenet", type: "dhcp"
    d.vm.network :forwarded_port, host: 3023, guest: 22
    d.vm.network :forwarded_port, host: 3080, guest: 80
    config.vm.provider :virtualbox do |vb|
      vb.customize ["modifyvm", :id, "--memory", "2048", "--cpus", "2", "--ioapic", "on"]
    end
    config.vm.provision "shell", inline: <<-SHELL
      route del default gw 10.0.2.2 eth0
      route add default gw 192.168.34.1 eth2
    SHELL
  end

end
