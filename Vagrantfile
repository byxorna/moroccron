# -*- mode: ruby -*-
# vi: set ft=ruby :

# All Vagrant configuration is done below. The "2" in Vagrant.configure
# configures the configuration version (we support older styles for
# backwards compatibility). Please don't change it unless you know what
# you're doing.

IP_ADDR = "10.10.0.5"
Vagrant.configure(2) do |config|
  # The most common configuration options are documented and commented below.
  # For a complete reference, please see the online documentation at
  # https://docs.vagrantup.com.

  # Every Vagrant development environment requires a box. You can search for
  # boxes at https://atlas.hashicorp.com/search.
  config.vm.box = "ubuntu/trusty64"

  config.vm.provider :virtualbox do |vb|
    vb.name = "moroccron-vagrant"
    vb.memory = "2048"
    vb.cpus = "2"
    #vb.customize ['modifyvm', :id, '--memory', "2048"]
    #vb.customize ['modifyvm', :id, '--cpus',   "2"]
  end

  config.vm.network "private_network", ip: IP_ADDR
  config.vm.network "forwarded_port", guest: 5050, host: 5050
  config.vm.network "forwarded_port", guest: 8000, host: 8000

  #config.vm.provision "file", source: "vagrant/profile", destination: "/home/vagrant/.profile"
  # be nice, so people can hack inside vbox
  config.vm.provision "file", source: "~/.gitconfig", destination: ".gitconfig"

  [
    ["1-system"],
    ["2-golang"],
  ].each do |s|
    config.vm.provision(:shell, privileged: false) do |shell|
      shell.path = "vagrant/#{s.first}"
      shell.args = s[1,s.length]
    end
  end

  config.vm.provision(:shell) do |shell|
    shell.path = "vagrant/3-mesosflexinstall"
    shell.args = ['--slave-hostname', IP_ADDR]
  end

  [
    ["4-services",IP_ADDR],
  ].each do |s|
    config.vm.provision(:shell, privileged: false) do |shell|
      shell.path = "vagrant/#{s.first}"
      shell.args = s[1,s.length]
    end
  end

  config.vm.synced_folder ".", "/home/vagrant/go/src/github.com/byxorna/moroccron"

  # Disable automatic box update checking. If you disable this, then
  # boxes will only be checked for updates when the user runs
  # `vagrant box outdated`. This is not recommended.
  # config.vm.box_check_update = false

  # Create a forwarded port mapping which allows access to a specific port
  # within the machine from a port on the host machine. In the example below,
  # accessing "localhost:8080" will access port 80 on the guest machine.
  # config.vm.network "forwarded_port", guest: 80, host: 8080

  # Create a private network, which allows host-only access to the machine
  # using a specific IP.
  # config.vm.network "private_network", ip: "192.168.33.10"

  # Create a public network, which generally matched to bridged network.
  # Bridged networks make the machine appear as another physical device on
  # your network.
  # config.vm.network "public_network"

  # Share an additional folder to the guest VM. The first argument is
  # the path on the host to the actual folder. The second argument is
  # the path on the guest to mount the folder. And the optional third
  # argument is a set of non-required options.
  # config.vm.synced_folder "../data", "/vagrant_data"

  # Provider-specific configuration so you can fine-tune various
  # backing providers for Vagrant. These expose provider-specific options.
  # Example for VirtualBox:
  #
  # config.vm.provider "virtualbox" do |vb|
  #   # Display the VirtualBox GUI when booting the machine
  #   vb.gui = true
  #
  #   # Customize the amount of memory on the VM:
  #   vb.memory = "1024"
  # end
  #
  # View the documentation for the provider you are using for more
  # information on available options.

  # Define a Vagrant Push strategy for pushing to Atlas. Other push strategies
  # such as FTP and Heroku are also available. See the documentation at
  # https://docs.vagrantup.com/v2/push/atlas.html for more information.
  # config.push.define "atlas" do |push|
  #   push.app = "YOUR_ATLAS_USERNAME/YOUR_APPLICATION_NAME"
  # end

  # Enable provisioning with a shell script. Additional provisioners such as
  # Puppet, Chef, Ansible, Salt, and Docker are also available. Please see the
  # documentation for more information about their specific syntax and use.
  # config.vm.provision "shell", inline: <<-SHELL
  #   sudo apt-get update
  #   sudo apt-get install -y apache2
  # SHELL
end
