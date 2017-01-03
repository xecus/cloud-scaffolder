# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|
  {{range $i1, $vm := .Vms}}
    config.vm.define "{{$vm.Hostname}}" do |d|
      d.vm.hostname = "{{$vm.Hostname}}"
      d.vm.box = "{{$vm.Image}}"
      {{range $i2, $ni := $vm.NetworkInterfaces}}d.vm.network {{$ni.Display}}
      {{end}}
      config.vm.provider :virtualbox do |vb|
        vb.customize ["modifyvm", :id, "--memory", "{{$vm.Memory}}", "--cpus", "{{$vm.Cpus}}", "--ioapic", "on"]
      end
      config.vm.provision "shell", inline: <<-SHELL
      SHELL
    end
  {{end}}
end
