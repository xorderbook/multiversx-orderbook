npm run build-devnet
zip -r build.zip build
ssh root@142.171.139.187 "rm -rf /var/www/html/devnet/*;"
scp build.zip root@142.171.139.187:/var/www/html/devnet
ssh root@142.171.139.187 "unzip -d /var/www/html/devnet /var/www/html/devnet/build.zip"
ssh root@142.171.139.187 "mv /var/www/html/devnet/build/* /var/www/html/devnet"