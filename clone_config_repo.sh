if [ -d "config-repo" ]; then 
    rm config-repo -rf
fi
if [ -d "config" ]; then 
    rm config -rf; 
fi
git clone git@github.com:SecureApiGateway/fr-platform-config.git config-repo
ls
configgithubsha=$(cat .config-git-ref)
cd config-repo 
git checkout $configgithubsha
cd ../
mkdir config
mv config-repo/config/* config
rm -rf config-repo