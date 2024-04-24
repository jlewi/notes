from peft import AutoPeftModelForCausalLM
from transformers import AutoTokenizer
model_id='hamel/hc-mistral-qlora-6'
model = AutoPeftModelForCausalLM.from_pretrained(model_id).cuda()
tokenizer = AutoTokenizer.from_pretrained(model_id)
tokenizer.pad_token = tokenizer.eos_token

def prompt(nlq, cols):
    return f"""[INST] <<SYS>>
Honeycomb AI suggests queries based on user input and candidate columns.
<</SYS>>

User Input: {nlq}

Candidate Columns: {cols}
[/INST]
"""

def prompt_tok(nlq, cols):
    _p = prompt(nlq, cols)
    input_ids = tokenizer(_p, return_tensors="pt", truncation=True).input_ids.cuda()
    out_ids = model.generate(input_ids=input_ids, max_new_tokens=5000, 
                          do_sample=False)
    return tokenizer.batch_decode(out_ids.detach().cpu().numpy(), 
                                  skip_special_tokens=True)[0][len(_p):]


if __name__ == "__main__":
    nlq = "Exception count by exception and caller"
    cols = ['error', 'exception.message', 'exception.type', 'exception.stacktrace', 'SampleRate', 'name', 'db.user', 'type', 'duration_ms', 'db.name', 'service.name', 'http.method', 'db.system', 'status_code', 'db.operation', 'library.name', 'process.pid', 'net.transport', 'messaging.system', 'rpc.system', 'http.target', 'db.statement', 'library.version', 'status_message', 'parent_name', 'aws.region', 'process.command', 'rpc.method', 'span.kind', 'serializer.name', 'net.peer.name', 'rpc.service', 'http.scheme', 'process.runtime.name', 'serializer.format', 'serializer.renderer', 'net.peer.port', 'process.runtime.version', 'http.status_code', 'telemetry.sdk.language', 'trace.parent_id', 'process.runtime.description', 'span.num_events', 'messaging.destination', 'net.peer.ip', 'trace.trace_id', 'telemetry.instrumentation_library', 'trace.span_id', 'span.num_links', 'meta.signal_type', 'http.route']

    out = prompt_tok(nlq, cols)
    print(out)