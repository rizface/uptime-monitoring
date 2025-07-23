# Chown the binary to current user
sudo chown $USER:$USER uptime-monitoring

exitCode=$?;

if [ $exitCode -ne 0 ]; then
    echo "Failed chown to $USER:$USER"
    exit $exitCode
fi

# Make the binary only executable by the owner
chmod 700 uptime-monitoring

# Move the binary to /usr/local/bin
sudo mv uptime-monitoring /usr/local/bin

# Create a directory for the database
mkdir -p /var/lib/uptime-monitoring 2> /dev/null

# Run
/usr/local/bin/uptime-monitoring /var/lib/uptime-monitoring/uptime.db