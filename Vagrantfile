# -*- mode: ruby -*-
# vi: set ft=ruby :

# All Vagrant configuration is done below. The "2" in Vagrant.configure
# configures the configuration version (we support older styles for
# backwards compatibility). Please don't change it unless you know what
# you're doing.
Vagrant.configure("2") do |config|
  # The most common configuration options are documented and commented below.
  # For a complete reference, please see the online documentation at
  # https://docs.vagrantup.com.

  # Every Vagrant development environment requires a box. You can search for
  # boxes at https://vagrantcloud.com/search.
  config.vm.box = "debian/stretch64"

  # Disable automatic box update checking. If you disable this, then
  # boxes will only be checked for updates when the user runs
  # `vagrant box outdated`. This is not recommended.
  # config.vm.box_check_update = false

  # Create a forwarded port mapping which allows access to a specific port
  # within the machine from a port on the host machine. In the example below,
  # accessing "localhost:8080" will access port 80 on the guest machine.
  # NOTE: This will enable public access to the opened port
  # config.vm.network "forwarded_port", guest: 80, host: 8080

  # Create a forwarded port mapping which allows access to a specific port
  # within the machine from a port on the host machine and only allow access
  # via 127.0.0.1 to disable public access
  # config.vm.network "forwarded_port", guest: 80, host: 8080, host_ip: "127.0.0.1"

  # Create a private network, which allows host-only access to the machine
  # using a specific IP.
  # config.vm.network "private_network", ip: "192.168.33.10"
  config.vm.network "private_network", ip: "192.168.100.100"

  # Create a public network, which generally matched to bridged network.
  # Bridged networks make the machine appear as another physical device on
  # your network.
  # config.vm.network "public_network"

  # Share an additional folder to the guest VM. The first argument is
  # the path on the host to the actual folder. The second argument is
  # the path on the guest to mount the folder. And the optional third
  # argument is a set of non-required options.
  # config.vm.synced_folder "../data", "/vagrant_data"
  config.vm.synced_folder ".", "/vagrant", type: "rsync",
    rsync__exclude: ".git/"

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

  # Enable provisioning with a shell script. Additional provisioners such as
  # Puppet, Chef, Ansible, Salt, and Docker are also available. Please see the
  # documentation for more information about their specific syntax and use.
  # config.vm.provision "shell", inline: <<-SHELL
  #   apt-get update
  #   apt-get install -y apache2
  # SHELL
  config.vm.provision "shell", inline: <<-SHELL
    apt-get update -qq
    apt-get install -qq -y dnsutils bind9 nginx curl vim git gcc

    curl -s -O -L https://dl.google.com/go/go1.11.1.linux-amd64.tar.gz
    tar -C /usr/local -xzf go1.11.1.linux-amd64.tar.gz
    rm go1.11.1.linux-amd64.tar.gz
    echo "export PATH=\\$PATH:/usr/local/go/bin" >> /home/vagrant/.bashrc

    curl -s -O -L https://nodejs.org/dist/v8.12.0/node-v8.12.0-linux-x64.tar.xz
    tar -C /usr/local -xJf node-v8.12.0-linux-x64.tar.xz
    rm node-v8.12.0-linux-x64.tar.xz
    mv /usr/local/node-* /usr/local/node
    echo "export PATH=\\$PATH:/usr/local/node/bin" >> /home/vagrant/.bashrc

    cat > /etc/nginx/sites-available/default <<EOF
    server {
      listen 80 default_server;
      listen [::]:80 default_server;
      root /vagrant/webui/dist;

      index index.html;
      server_name _;
      location / {
        try_files $uri $uri/ =404;
      }
      location /api {
            return 302 /api/;
      }
      location /api {
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_pass http://localhost:8000;
      }
    }
    EOF
    systemctl reload nginx

    source /home/vagrant/.bashrc
    su - vagrant -c "cd /vagrant/api && go get ."
    su - vagrant -c "cd /vagrant/webui && npm install"

    cat > /etc/systemd/system/koala.service << EOF
    [Unit]
    Description=Koala DNS editing frontend
    After=network.target

    [Service]
    Type=simple
    User=root
    Environment=KOALA_ADDR=localhost:8000
    Environment=KOALA_ZONEFILE=/vagrant/api/example.zone
    Environment=KOALA_APPLYCMD="systemctl reload bind9"
    WorkingDirectory=/root/
    ExecStart=/home/vagrant/go/bin/api
    Restart=on-abort

    [Install]
    WantedBy=multi-user.target
    EOF
    systemctl daemon-reload && systemctl enable koala && systemctl start koala
  SHELL
end
