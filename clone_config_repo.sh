# Clone line can be used in future if versioning is put in place
# git clone --depth 1 --branch ${github-tag} git@github.com:SecureApiGateway/fr-platform-config.git config-repo

#Script for cloning a version of the fr-platform-config repo, version to clone is saved in .config-git-ref. Files are then moved to a desired folder


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