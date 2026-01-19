# Suggested Ollama Models for Tact

This document provides recommendations for Ollama models suitable for Tact's time entry parsing system. The system parses natural language time entries, extracts duration/time codes/work types, outputs structured JSON, and performs semantic matching with RAG context.

## Top 5 Recommendations

### 1. Qwen2.5:7b (Recommended)

| Pros | Cons |
|------|------|
| Excellent instruction following and JSON output | Slightly larger than current 3b model |
| Strong reasoning for semantic matching | Less community testing than Llama |
| Native function calling support in Ollama | |
| Great multilingual support if you have international time codes | |
| Active development from Alibaba with frequent updates | |

**Why for Tact**: The strong reasoning capabilities help with semantic matching ("APHL standup" → FEDS-163), and Qwen models are specifically praised for structured output reliability.

---

### 2. Llama3.2:3b (Current Default)

| Pros | Cons |
|------|------|
| Very lightweight (4GB RAM minimum) | Smaller context can miss nuances |
| Fast inference for background tasks | Less sophisticated reasoning than 7b+ models |
| Native Ollama function calling support | May struggle with ambiguous entries |
| Well-tested, large community | |

**Why for Tact**: Good balance for a background parsing task. If accuracy is acceptable today, this remains a solid choice. Consider upgrading only if you see parsing quality issues.

---

### 3. Hermes3:8b

| Pros | Cons |
|------|------|
| 90% accuracy on function calling benchmarks | Not an official Ollama library model (community) |
| Purpose-built for structured JSON output | Slightly less general-purpose |
| Excellent at following complex system prompts | 8b requires more RAM (~8GB min) |
| Specifically trained on JSON mode datasets | |

**Why for Tact**: If you're seeing JSON parsing errors or malformed outputs, Hermes3 is specifically designed for this exact use case. It has "powerful and reliable function calling and structured output capabilities."

---

### 4. Mistral-Small:22b (or Mistral:7b for lighter option)

| Pros | Cons |
|------|------|
| "Low-latency function calling" optimization | 22b version needs significant RAM (16GB+) |
| Strong instruction following | 7b version less capable than Qwen2.5:7b |
| Fast inference relative to size | |
| Good at concise, focused responses | |

**Why for Tact**: Mistral Small 3 was specifically updated for "improved function calling, instruction following, and less repetition errors." Good if you want a balance of speed and capability.

---

### 5. Llama3.1:8b (or Llama3.3:70b for max quality)

| Pros | Cons |
|------|------|
| 128K context window (great for large RAG contexts) | 8b is middle-ground, not best at anything |
| Most widely deployed/tested | 70b requires serious hardware (32GB+ RAM) |
| Excellent general reasoning | |
| Strong ecosystem support | |

**Why for Tact**: The 128K context is useful if you have many time codes/rules in RAG context. Llama 3.1 is described as a "workhorse" that "adapts to tasks from email drafting to data summarization."

---

## Quick Decision Guide

| Scenario | Recommended Model |
|----------|-------------------|
| Current accuracy is acceptable | Stick with `llama3.2:3b` |
| Want better accuracy without much more resources | Try `qwen2.5:7b` |
| Getting JSON formatting issues | Try `hermes3:8b` |
| Have lots of RAM and want best quality | Try `llama3.3:70b` |
| Need fast inference with good quality | Try `mistral:7b` |

## Configuration

Set the model via the `TACT_OLLAMA_MODEL` environment variable:

```bash
# Examples
TACT_OLLAMA_MODEL=llama3.2:3b      # Default, lightweight
TACT_OLLAMA_MODEL=qwen2.5:7b       # Recommended upgrade
TACT_OLLAMA_MODEL=hermes3:8b       # Best for JSON reliability
TACT_OLLAMA_MODEL=mistral:7b       # Fast with good quality
```

## Hardware Requirements

| Model Size | Minimum RAM | Recommended RAM |
|------------|-------------|-----------------|
| 1B-3B | 4GB | 8GB |
| 7B-8B | 8GB | 16GB |
| 13B-14B | 16GB | 32GB |
| 22B+ | 32GB | 64GB |
| 70B+ | 64GB | 128GB |

## References

- [Ollama Structured Outputs](https://docs.ollama.com/capabilities/structured-outputs)
- [Ollama Models Comparison 2025](https://collabnix.com/best-ollama-models-in-2025-complete-performance-comparison/)
- [Hermes 3 on Ollama](https://ollama.com/library/hermes3)
- [Mistral Small on Ollama](https://ollama.com/library/mistral-small)
- [Open Source LLMs 2025](https://huggingface.co/blog/daya-shankar/open-source-llms)
- [Choosing Ollama Models Guide](https://collabnix.com/choosing-ollama-models-the-complete-2025-guide-for-developers-and-enterprises/)
