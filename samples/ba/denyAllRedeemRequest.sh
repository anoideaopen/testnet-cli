#!/bin/bash

set -xe

echo "=========================================================================="
echo "=========================================================================="
echo "=========================================================================="
./cli query ba redeemRequestsList --connection atmz-ch-preprod/config_test.yaml
./cli query gf28iln060f redeemRequestsList --connection atmz-ch-preprod/config_test.yaml
./cli query gf28iln060 redeemRequestsList --connection atmz-ch-preprod/config_test.yaml
./cli query gpfbar047053 redeemRequestsList --connection atmz-ch-preprod/config_test.yaml
./cli query gpfbar079001 redeemRequestsList --connection atmz-ch-preprod/config_test.yaml
./cli query etc redeemRequestsList --connection atmz-ch-preprod/config_test.yaml
./cli query atz029olp005xx redeemRequestsList --connection atmz-ch-preprod/config_test.yaml
./cli query etc1 redeemRequestsList --connection atmz-ch-preprod/config_test.yaml
./cli query vt requestsList --connection atmz-ch-preprod/config_test.yaml
./cli query it redeemRequestsList --connection atmz-ch-preprod/config_test.yaml
echo "=========================================================================="
echo "=========================================================================="
echo "=========================================================================="
#./cli invoke ba denyAllRedeemRequest -s 3aDebSkgXq37VPrzThboaV8oMMbYXrRAt7hnGrod4PNMnGfXjh14TY7cQs8eVT46C4RK4ZyNKLrBmyD5CYZiFmkr --connection atmz-ch-preprod/config_test.yaml
#./cli invoke gf28iln060f denyAllRedeemRequest -s 3aDebSkgXq37VPrzThboaV8oMMbYXrRAt7hnGrod4PNMnGfXjh14TY7cQs8eVT46C4RK4ZyNKLrBmyD5CYZiFmkr --connection atmz-ch-preprod/config_test.yaml
#./cli invoke gf28iln060 denyAllRedeemRequest -s 3aDebSkgXq37VPrzThboaV8oMMbYXrRAt7hnGrod4PNMnGfXjh14TY7cQs8eVT46C4RK4ZyNKLrBmyD5CYZiFmkr --connection atmz-ch-preprod/config_test.yaml
#./cli invoke gpfbar047053 denyAllRedeemRequest -s 3aDebSkgXq37VPrzThboaV8oMMbYXrRAt7hnGrod4PNMnGfXjh14TY7cQs8eVT46C4RK4ZyNKLrBmyD5CYZiFmkr --connection atmz-ch-preprod/config_test.yaml
#./cli invoke gpfbar079001 denyAllRedeemRequest -s 3aDebSkgXq37VPrzThboaV8oMMbYXrRAt7hnGrod4PNMnGfXjh14TY7cQs8eVT46C4RK4ZyNKLrBmyD5CYZiFmkr --connection atmz-ch-preprod/config_test.yaml
#./cli invoke etc denyAllRedeemRequest -s 3aDebSkgXq37VPrzThboaV8oMMbYXrRAt7hnGrod4PNMnGfXjh14TY7cQs8eVT46C4RK4ZyNKLrBmyD5CYZiFmkr --connection atmz-ch-preprod/config_test.yaml
#./cli invoke atz029olp005xx denyAllRedeemRequest -s 3aDebSkgXq37VPrzThboaV8oMMbYXrRAt7hnGrod4PNMnGfXjh14TY7cQs8eVT46C4RK4ZyNKLrBmyD5CYZiFmkr --connection atmz-ch-preprod/config_test.yaml
#./cli invoke etc1 denyAllRedeemRequest -s 3aDebSkgXq37VPrzThboaV8oMMbYXrRAt7hnGrod4PNMnGfXjh14TY7cQs8eVT46C4RK4ZyNKLrBmyD5CYZiFmkr --connection atmz-ch-preprod/config_test.yaml
#./cli invoke vt denyAllDistribRequest -s 3aDebSkgXq37VPrzThboaV8oMMbYXrRAt7hnGrod4PNMnGfXjh14TY7cQs8eVT46C4RK4ZyNKLrBmyD5CYZiFmkr --connection atmz-ch-preprod/config_test.yaml
#./cli invoke it denyAllRedeemRequest -s 3aDebSkgXq37VPrzThboaV8oMMbYXrRAt7hnGrod4PNMnGfXjh14TY7cQs8eVT46C4RK4ZyNKLrBmyD5CYZiFmkr --connection atmz-ch-preprod/config_test.yaml
#./cli invoke it denyAllDistribRequest -s 3aDebSkgXq37VPrzThboaV8oMMbYXrRAt7hnGrod4PNMnGfXjh14TY7cQs8eVT46C4RK4ZyNKLrBmyD5CYZiFmkr --connection atmz-ch-preprod/config_test.yaml
echo "=========================================================================="
echo "=========================================================================="
echo "=========================================================================="
