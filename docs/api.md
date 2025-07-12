# API Documentation

## Core Commands

### `allora ask`
Ask AI questions about your infrastructure.

```bash
allora ask "How can I optimize my AWS costs?"
allora ask "What's causing high latency in my application?"
```

### `allora deploy`
Deploy applications and infrastructure.

```bash
allora deploy --environment production --service web-app
allora deploy kubernetes --app myapp --namespace production
```

### `allora monitor`
Monitor infrastructure and applications.

```bash
allora monitor --provider aws --resource ec2
allora monitor metrics --service database
```

### `allora troubleshoot`
AI-powered troubleshooting.

```bash
allora troubleshoot --service database --issue "high latency"
allora troubleshoot logs --application web-app
```

### `allora security`
Security analysis and compliance.

```bash
allora security audit --provider aws
allora security compliance --standard SOC2
```

## Configuration Commands

### `allora config`
Manage configuration.

```bash
allora config init
allora config set aws.region us-west-2
allora config get aws.region
allora config validate
```

## Plugin Commands

### `allora plugin`
Manage plugins.

```bash
allora plugin list
allora plugin install <plugin-name>
allora plugin remove <plugin-name>
```

## Interactive Mode

### `allora gemini`
Launch the interactive Gemini-style interface.

```bash
allora gemini
allora gemini --export conversation.json
```

## Global Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--config` | Configuration file path | `~/.allora/config.yaml` |
| `--log-level` | Log level (debug, info, warn, error) | `info` |
| `--output` | Output format (json, yaml, table) | `table` |
| `--no-color` | Disable colored output | `false` |
