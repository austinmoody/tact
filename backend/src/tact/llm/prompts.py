from tact.llm.provider import ParseContext

SYSTEM_PROMPT_TEMPLATE = """You are a time entry parser. \
Your job is to extract structured information from natural language time entries.
{rag_context_text}
Available Time Codes:
{time_codes_text}

Available Work Types:
{work_types_text}

Instructions:
1. Extract the duration in minutes. Examples:
   - "10m" = 10 minutes (m = minutes)
   - "30 min" = 30 minutes
   - "1h" = 60 minutes
   - "2h" = 120 minutes
   - "1h30m" = 90 minutes
   - "1.5h" = 90 minutes
   IMPORTANT: "m" always means minutes, NOT hours. "10m" = 10, not 600.
2. Match to a time_code_id using this priority order:
   a. FIRST check Matching Context Rules above - if the entry text matches a rule, use that time_code_id
   b. THEN check time code keywords and descriptions
   Example: If context says "(time_code: ALM-123): Vibe coding" and entry mentions "vibe coding", use ALM-123
3. Match to a work_type_id from the available list
4. Generate a clean description of the work done
5. Provide confidence scores (0.0 to 1.0) for each field
6. Provide brief reasoning notes explaining your matching decision

Respond with ONLY valid JSON in this exact format:
{{
  "duration_minutes": <integer or null>,
  "time_code_id": "<string or null>",
  "work_type_id": "<string or null>",
  "parsed_description": "<string or null>",
  "confidence_duration": <float 0-1>,
  "confidence_time_code": <float 0-1>,
  "confidence_work_type": <float 0-1>,
  "confidence_overall": <float 0-1>,
  "notes": "<string explaining reasoning for matches, especially which context rules applied>"
}}

If you cannot determine a field, set it to null and give low confidence.
In notes, explain why you made each decision - especially mention which context rules influenced the match.
Do not include any text outside the JSON object."""


def build_system_prompt(context: ParseContext) -> str:
    """Build the system prompt with time codes, work types, and RAG context."""
    # Build RAG context section
    rag_context_text = ""
    if context.rag_contexts:
        rag_lines = ["\nMatching Context Rules (use these to assign time_code_id):"]
        for rc in context.rag_contexts:
            if rc.time_code_id:
                rag_lines.append(f"- If entry mentions \"{rc.content}\" → use time_code_id: {rc.time_code_id}")
            else:
                rag_lines.append(f"- Project {rc.project_id} rule: {rc.content}")
        rag_context_text = "\n".join(rag_lines) + "\n"

    time_codes_text = "\n".join(
        f"- {tc.id}: {tc.name} - {tc.description} "
        f"(keywords: {', '.join(tc.keywords)})"
        for tc in context.time_codes
    )
    if not time_codes_text:
        time_codes_text = "(none defined)"

    work_types_text = "\n".join(
        f"- {wt.id}: {wt.name}" for wt in context.work_types
    )
    if not work_types_text:
        work_types_text = "(none defined)"

    return SYSTEM_PROMPT_TEMPLATE.format(
        rag_context_text=rag_context_text,
        time_codes_text=time_codes_text,
        work_types_text=work_types_text,
    )


def build_user_prompt(user_input: str) -> str:
    """Build the user prompt with the entry to parse."""
    return f'Parse this time entry:\n"{user_input}"'
