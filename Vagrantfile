Vagrant.configure("2") do |config|
  config.vm.box = "centos/7"
  config.vm.network "private_network", type: "dhcp"

  config.vm.provider "virtualbox" do |v|
    v.customize ["modifyvm", :id, "--natdnshostresolver1", "on"]

    v.name = "fuse_builder"
    v.memory = 16384
    v.cpus = 6
  end

  config.vm.provision "shell", inline: <<-BASH
    # install fuse
    sudo yum install -y curl fuse* afuse git
    sudo modprobe fuse

    # install go
    if [ ! -f /usr/local/go/bin/go ]; then
        curl -O -L https://dl.google.com/go/go1.12.linux-amd64.tar.gz
        tar -C /usr/local -xzf go1.12.linux-amd64.tar.gz
        echo "export PATH=/usr/local/go/bin:/root/go/bin:\$PATH" >> /etc/profile
        echo "export GOROOT=/usr/local/go" >> /etc/profile
        . /etc/profile
        go version

        mkdir -p /root/go/src /root/go/bin /home/vagrant/go/src /home/vagrant/go/bin
    fi

    # install dep
    if [ ! -f /usr/local/go/bin/dep ]; then
        curl https://raw.githubusercontent.com/golang/dep/master/install.sh | INSTALL_DIRECTORY=/usr/local/go/bin sh
        dep version
    fi

#    rm -rf /home/vagrant/go/src/docs-fuse
#    rsync -avz \
#        --exclude '.vagrant' \
#        --exclude '.git' \
#        --exclude '.idea' \
#        --exclude 'vendor' \
#        /vagrant/ /home/vagrant/go/src/docs-fuse
#
#    cd /home/vagrant/go/src/docs-fuse
#    chown vagrant:vagrant -R /home/vagrant/go/
#    sudo su vagrant
#    . /etc/profile


    rm -rf /root/go/src/docs-fuse
    rsync -avz \
        --exclude '.vagrant' \
        --exclude '.git' \
        --exclude '.idea' \
        --exclude 'vendor' \
        /vagrant/ /root/go/src/docs-fuse


    cd /root/go/src/docs-fuse
    dep ensure
    go build -o docs.fuse cmd/docs-fuse-raw/main.go

    mkdir -p /mnt/sismics

    ./docs.fuse -d /mnt/sismics

BASH

end
