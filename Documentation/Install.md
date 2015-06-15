Setup
=====

## On a dev environment

### Theory

You must set an environment with Golang. In your GOPATH, create a directory supinfo/mewpipe and clone the project inside it.
Run ./fixtures.sh to get some example data.
Run ./test.sh to run the test suit on http://localhost:8080
Run ./run.sh to run the server
Run ./build.sh to cross compile the project.

### Practice

On a fresh debian 8.1 install with sudo :

    #Install Golang
    $ sudo apt-get install git
    $ wget https://storage.googleapis.com/golang/go1.4.2.linux-amd64.tar.gz
    $ sudo tar -C /usr/local -xzf go1.4.2.linux-amd64.tar.gz
    $ echo 'export PATH=$PATH:/usr/local/go/bin' | sudo tee -a /etc/profile
    $ echo 'export GOPATH=/home/<User>' | tee -a ~/.profile
    $ source /etc/profile
    $ source ~/.profile
    $ rm go1.4.2.linux-amd64.tar.gz
    $ mkdir -p src/supinfo
    $ cd src/supinfo 
    $ git clone git@spider4all.com:4PJT-MewPipe mewpipe #Available on Github after the presentation @ Supinfo
    $ mv 4PJT-MewPipe src/supinfo/mewpipe
    $ cd mewpipe/
    $ go get
    $ cp configs/local.ini.sample configs/local.ini
    #Override vars from configs/base.ini in configs/local.ini like mongo connection string, database or http binding
    $ ~/src/supinfo/mewpipe/run.sh

## Deploy to prod

### Build

Run ./build.sh on dev environment, like : 

    $ ./build.sh 
    Cross building, outputs are in /to/my/path/mewpipe/build
    Building for linux on 386...
    Compress build/mewpipe_linux_386.tar.bz2
    ---------------------------
    Building for linux on amd64...
    Compress build/mewpipe_linux_amd64.tar.bz2
    ---------------------------
    Building for linux on arm...
    Compress build/mewpipe_linux_arm.tar.bz2
    ---------------------------
    Building for windows on 386...
    Compress build/mewpipe_windows_386.tar.bz2
    ---------------------------
    Building for windows on amd64...
    Compress build/mewpipe_windows_amd64.tar.bz2
    ---------------------------
    Building for darwin on 386...
    Compress build/mewpipe_darwin_386.tar.bz2
    ---------------------------
    Building for darwin on amd64...
    Compress build/mewpipe_darwin_amd64.tar.bz2
    ---------------------------

### Upload to Prod

Choose the right build and upload it on your server or prod environment. Your right build is mewpipe_{Operating kernel}_{processor infrastructure}.tar.bz2
In my case, on my mac OS X it's mewpipe_darwin_amd64.tar.bz2, on my debian it's mewpipe_linux_amd64.tar.bz2, on my Raspberry Pi it's mewpipe_linux_arm.tar.bz2

On my dev environment, I use the following command :

    $ scp build/mewpipe_linux_amd64.tar.bz2 debian@<my ip>:mewpipe.tar.bz2


On my prod server, I use 

    $ tar -xjf mewpipe.tar.bz2
    $ rm mewpipe.tar.bz2
    $ chmod +x mewpipe
    $ mewpipe --help # This show all variables you can overload on config.ini or pass in parameter.
    $ ./mewpipe -config="config.ini" # You are free to deamonize this process.