# Sandwhich-Server

The backend REST API for Sandwhich

## Getting Started
* Install Golang (sample script below for Ubuntu 64bit)
* Run ./bin/getimports.sh
* go build && ./sandwhich
* for usage details: ./sandwhich -help

```bash
### Go installation script
# Remove previous go installations
sudo rm -r /usr/local/go
cd /tmp
sudo rm -r go

# Download go (Update the version/OS here)
wget https://storage.googleapis.com/golang/go1.3.3.linux-amd64.tar.gz
tar -zxf go1.3.3.linux-amd64.tar.gz

# Move go folder to default folder /usr/local
sudo mv go /usr/local/

# Add to path and .profile (comment out if updating)
export PATH=$PATH:/usr/local/go/bin
echo "export PATH=\$PATH:/usr/local/go/bin" >> $HOME/.profile
```