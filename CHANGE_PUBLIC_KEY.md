# The process of updating the public key

# 1. Message generation

Admin generates a signature message on linux

Parameters:1. Public keys of all validators that were specified during installation/initialization of acl.

It is mandatory to sign with 3 validators to update the public key.



```shell
export validatorPublicKeys="A4JdE9iZRzU9NEiVDNxYKKWymHeBxHR7mA8AetFrg8m4,\
5Tevazf8xxwyyKGku4VCCSVMDN56mU3mm2WsnENk1zv5,\
6qFz88dv2R8sXmyzWPjvzN6jafv7t1kNUHztYKjH1Rd4"
```

1. The address to change the public key, format base58sheck:
```shell
export changedAddr="2GFkmC1RE1kMe1HcdcW9Sk3d7nBgNtxeDPRXK7xxxrXordZa7b"
```
1. New public key, base58 format:
```shell
export newPkey="BREP5CVURcJ6CoTdUpJdSNgZThMXvubFueyviSHDRW4Y"
```
Example of running the command for ACL 0.2.0:
```shell
./cli generateMessage acl changePublicKey $validatorPublicKeys $changedAddr $newPkey
```

Example of running the command for ACL 0.3.1:

**reason** - Comment text field. For example: "On the basis of writ of execution #blah-blah-blah-blah of such-and-such date"
**reasonId** - integer field, which will be defined by enum on the backend (for example: 0-transfer of wallet control by law, 1-return of control to user, 2-loss of user's public key).

```shell
export reason="lost_key"
export reasonId=2
./cli generateMessage acl changePublicKey $validatorPublicKeys $changedAddr $reason $reasonId $newPkey
``` 

Result:
A message.txt file appears in the directory next to the "cli" application

The archive for the validator must contain:
- message.txt
- cli-windows-amd64.exe
- cli-windows-386.exe
  Create and transfer a zip archive to Validator.

# 2. Message Signature

The validator must run the command line on Windows:
1. Double-click on the run_cmd.bat file
2. Replace SECRET_KEY_PUT_HERE with a private key in base58 or hex format.
   COMMAND TEMPLATE:

```shell
cli-windows-amd64.exe -s SECRET_KEY_PUT_HERE signMessage
```

An example of the command you should end up with:
```shell
cli-windows-amd64.exe -s 3aDebSkgXq37VPrzThboaV8oMMbYXrRAt7hnGrod4PNMnGfXjh14TY7cQs8eVT46C4RK4ZyNKLrBmyD5CYZiFmkr signMessage
```
Attention depending on the windows system you need to run either cli-windows-amd64.exe or cli-windows-386.exe

3. Enter the command in the open conque (after step 1) and press ENTER to execute the command.

4. As a result of executing the command, a file signature-*.txt appeared in the directory next to the "cli" application.
   The public key of the person who signed the message is specified in the file name instead of *.
   Copy this file and send it to the Administrator.

# 3. Preparing to update the public key in hlf
3.1. The administrator collects signature-*.txt files from validators and puts them in a directory next to the "cli" application.
3.2. In this step, you need to check that there is a message next to the cli application message.txt
3.3. The administrator configures the config_test.yaml file, and saves the crypto materials for hlf to the crypto folder as described in the config_test.yaml configuration file
3.3.1. Before you start, you need to save the config_test_tmp.yaml file named config_test.yaml. The config_test.yaml file is used by default when calling cli, as a configuration file for hlf. If you want to use a file with a different name, it can be specified with the `./cli --cfg config_test.yaml ....` or `./cli -f config_test.yaml ....` parameter.
3.3.2. Modify the configuration file for hlf config_test.yaml.
3.3.3. Channel name and policy. In the file you need to set the name of the channel where we make a request, as well as the list of peers in accordance with the policy for your channel.In my case, the name of the channel acl. 22 line in config_test_tmp.yaml file. **Attention domains should correspond to your stand.

```yaml
acl:
    peers:
        peer0.testnet.uat.dlt.testnet.ch: { }
        peer1.testnet.uat.dlt.testnet.ch: { }
        peer0.trafigura.uat.dlt.testnet.ch: { }
        peer0.traxys.uat.dlt.testnet.ch: { }
        peer0.umicore.uat.dlt.testnet.ch: { }
```

All of these peers must also be specified in the config_test.yaml configuration file in the peers block 56 line config_test_tmp.yaml
Let's specify in the config_test.yaml file the organizations that are in our hlf network. Line 23 in the config_test_tmp.yaml file.
Further it is necessary to copy the name of the organization from which we will execute requests and specify the name of the organization in the file config_test.yaml on line 5.
In my example this organization is called **testnet** organization: testnet
3.3.4. check that all paths to files and directories specified in config_test.yaml file are correct (example of crypto folder in the attachment).
3.3.5. config_test_tmp.yaml - you can delete it, we don't need it anymore because the settings are specified in the file config_test.yaml

# 4. Updating the public key in hlf

The administrator sends a request to hlf to change the public key. The user for hlf can be specified using the -u parameter

```shell
./cli sendRequest acl changePublicKey $validatorPublicKeys -u User14
```

Checking the public key in hlf
The administrator checks whether the public key for an address has changed or not by executing the following commands.

```shell
OLD_PublicKey="CTmpLBcWAtikpFYDwPkSPeQpKZALxpaGG7r5AYMiCjbG"
./cli query acl checkKeys $OLD_PublicKey -u User1

NEW_PublicKey="BREP5CVURcJ6CoTdUpJdSNgZThMXvubFueyviSHDRW4Y"
./cli query acl checkKeys $NEW_PublicKey -u User1
```