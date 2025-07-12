# Troubleshooting Guide

## Common Issues

### Installation Issues

#### Go Version Compatibility
**Error**: `go: module requires Go 1.21 or later`
**Solution**: Update Go to version 1.21 or higher.

#### Build Failures
**Error**: Build fails with compilation errors
**Solution**: 
1. Run `go mod tidy` to clean dependencies
2. Ensure you're in the correct directory
3. Check Go version compatibility

### Configuration Issues

#### API Key Not Working
**Error**: `authentication failed`
**Solution**:
1. Verify API key is correct
2. Check environment variables
3. Ensure proper permissions

#### Cloud Provider Connection Failed
**Error**: `failed to connect to AWS/Azure/GCP`
**Solution**:
1. Verify credentials
2. Check network connectivity
3. Validate region settings

### Runtime Issues

#### High Memory Usage
**Symptoms**: Application consuming excessive memory
**Solution**:
1. Reduce batch sizes
2. Enable streaming mode
3. Check for memory leaks in plugins

#### Slow Response Times
**Symptoms**: Commands taking too long
**Solution**:
1. Check network latency
2. Optimize queries
3. Use caching

## Debug Mode

Enable debug logging:

```bash
allora --log-level debug <command>
```

Or set environment variable:
```bash
export ALLORA_LOG_LEVEL=debug
```

## Getting Help

1. Check the [FAQ](faq.md)
2. Search [GitHub Issues](https://github.com/AlloraAi/AlloraCLI/issues)
3. Join our [Discord](https://discord.gg/alloracli)
4. Email support@alloracli.com
