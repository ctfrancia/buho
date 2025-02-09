# buho
API for Buho

## Description
Búho is a platform that allows you to create and manage your chess tournaments.
It is a tool that helps you to organize your tournaments, manage your players and games, and keep track of the results.
Búho is designed to be easy to use and flexible, so you can adapt it to your needs.


## Running the project
To run this there are some things that you wwill need to do:
1. private/public rsa keys (this is pem files and not)
  - You can generate them with the following command:
  ```bash
    openssl genrsa -out test_key.pem 2048
    # Don't add passphrase
    openssl rsa -in test_key.pem -outform PEM -pubout -out test_key.pem.pub
    ```
    - You will need to add the public key to the `keys` folder
    - (here see the article)[https://www.digitalocean.com/community/tutorials/openssl-essentials-working-with-ssl-certificates-private-keys-and-csrs]
    - or [here](https://auth0.com/docs/secure/application-credentials/generate-rsa-key-pair)
