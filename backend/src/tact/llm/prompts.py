from tact.llm.provider import ParseContext

SYSTEM_PROMPT_TEMPLATE = """You are a time entry parser. \
Your job is to extract structured information from natural language time entries.
{rag_context_text}
Available Time Codes:
{time_codes_text}

Available Work Types:
{work_types_text}

Instructions:
1. Extract the duration in minutes (e.g., "2h" = 120, "30 min" = 30)
2. Match to a time_code_id from the available list based on keywords and context
3. Match to a work_type_id from the available list
4. Generate a clean description of the work done
5. Provide confidence scores (0.0 to 1.0) for each field
6. Provide brief reasoning notes explaining your matching decision
7. IMPORTANT: If matching context rules are provided above, follow them carefully - \
they contain project-specific categorization rules that override generic matching

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
        rag_lines = ["\nMatching Context Rules:"]
        for rc in context.rag_contexts:
            if rc.time_code_id:
                source = f"(time_code: {rc.time_code_id})"
            else:
                source = f"(project: {rc.project_id})"
            rag_lines.append(f"- {source}: {rc.content}")
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
