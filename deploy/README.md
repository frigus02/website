# Deployment

This site is deployed on Uberspace.

## Setup

1. Follow the [Uberspace Node.js](https://wiki.uberspace.de/development:nodejs) documentation to configure `node`, `npm` and install `yarn`.

1. Clone this repository:

   ```sh
   mkdir -p ~/projects
   cd ~/projects
   git clone https://github.com/frigus02/website.git
   ```

1. Create script to handle GitHub webhooks:

   ```sh
   ~/projects/website/deploy/setup.sh
   ```

1. Setup a webhook in GitHub with the following settings:

   | Name        | Value                                            |
   | ----------- | ------------------------------------------------ |
   | Payload URL | `https://<domain><path-printed-by-prev-command>` |
   | Secret      | _empty_                                          |
   | Events      | Just the `push` event                            |
