# Shakespeare

Shakespeare is a tracer generator for the tempest benchmarking tool

## Examples

### Creating an OpenAI Benchmarking trace

```bash
shakespeare -o ./temp/test.json -t openai -g '{ 
        "rps": 1, 
        "duration": 10, 
        "apikey": "",
        "endpoint": "v1/chat/completions"
        "model": "modelName",
        "maxtokens": 100,
    }'
```