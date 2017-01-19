# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|
  {{range $i1, $vm := .Vms}}
    config.vm.define "{{$vm.Hostname}}" do |d|
      d.vm.hostname = "{{$vm.Hostname}}"
      d.vm.box = "{{$vm.Image.ImageName}}"
      {{range $i2, $ni := $vm.NetworkInterfaces}}d.vm.network {{$ni.ExpandNetworkInterfaceOptions}}
      {{end}}
      config.vm.provider :virtualbox do |vb|
        vb.customize ["modifyvm", :id, "--memory", "{{$vm.MemorySize}}", "--cpus", "{{$vm.NumOfCpus}}", "--ioapic", "on"]
      end
      config.vm.provision "shell", inline: <<-SHELL
      SHELL
    end
  {{end}}
end
