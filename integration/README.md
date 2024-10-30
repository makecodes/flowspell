# About the integration directory

This directoy is a manual integration tests and is useful to understand how Flowspell can be integrated with other systems.

## Commands

There are some commands which should be executed

### prepare

This the #1 command to be executed. It will prepare the data creating flows and task definitions

```bash
make prepare
```

### start

This command will start a flow instance which was previously created by the prepare command

```bash
make start
```
