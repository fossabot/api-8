
chmod 600 .ssh/id_rsa
apk add --no-cache curl rsync openssh
timestamp=$(date +%s)
ssh -i .ssh/id_rsa -o StrictHostKeyChecking=no bejo@$STAGING_HOST "mkdir -p ~/api/releases/$timestamp"
rsync -avz --exclude '.git' -e 'ssh -i .ssh/id_rsa -o StrictHostKeyChecking=no' ./api/out/api bejo@$STAGING_HOST:~/api/releases/$timestamp
ssh -i .ssh/id_rsa -o StrictHostKeyChecking=no bejo@$STAGING_HOST <<EOF
    set -x
    rm -f ~/api/current
    ln -s ~/api/releases/$timestamp ~/api/current
    sudo systemctl restart api.service
    cd ~/api/releases
    ls -t | tail -n +3 | xargs rm -rf --
EOF
