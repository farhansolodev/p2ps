#!/bin/bash -i

apt update
apt install -y unzip net-tools

# asdf
git clone https://github.com/asdf-vm/asdf.git ~/.asdf --branch v0.15.0
cat <<EOF >> ~/.bashrc
. "$HOME/.asdf/asdf.sh"
. "$HOME/.asdf/completions/asdf.bash"
EOF
. ~/.bashrc

# go
asdf plugin add golang https://github.com/asdf-community/asdf-golang.git
asdf install golang latest
asdf global golang latest

#stun
# go install github.com/pion/stun/v3/cmd/stun-nat-behaviour@latest
# cp /root/.asdf/installs/golang/1.23.4/packages/bin/stun-nat-behaviour ~
# ./stun-nat-behaviour

#p2p
git clone https://github.com/farhansolodev/p2pc
cd p2pc
go build .
