import os
from dotenv import load_dotenv
from langchain.agents import create_openai_tools_agent, AgentExecutor
from langchain_openai import ChatOpenAI
from langchain_core.prompts import ChatPromptTemplate, MessagesPlaceholder
from langchain.tools import tool
import argparse

# Load environment variables from .env file
load_dotenv()
os.environ["OPENAI_API_KEY"] = os.getenv("OPENAI_API_KEY")

# --- Define the Tools ---
# These tools allow the LLM to 'read' the simulated job files.

@tool
def read_file(file_path: str) -> str:
    """
    Reads the contents of a file at the given file path and returns it as a string. Returns an error message if the file cannot be read.
    """
    try:
        with open(file_path, "r") as f:
            content = f.read()
    except FileNotFoundError:
        return f"Error: File not found at {file_path}."
    except Exception as e:
        return f"Error reading file {file_path}: {e}"
    return content

# --- The Agent Logic ---

def analyze_job_failure(job_id: str, job_script_path: str, stdout_path: str, stderr_path: str) -> str:
    """
    Analyzes a simulated failed PBS Pro job and generates a diagnostic report.
    """

    # Create the LLM instance
    llm = ChatOpenAI(model="gpt-4-turbo-preview", temperature=0.0)

    # Define the tools available to the agent
    tools = [read_file]

    # The prompt instructs the LLM on its role and how to use the tools.
    # It emphasizes the order of reading files and the desired output format.
    prompt = ChatPromptTemplate.from_messages([
        ("system", f"""You are an expert HPC system administrator for a PBS Pro cluster.
        Your primary task is to diagnose why a user's batch job (Job ID: {job_id}) failed and provide a clear, actionable explanation.
        You have access to the user's job script, their standard output, and standard error files.
        You MUST use the `read_file` tool to examine the contents of these files to gather information.
        DO NOT invent file contents or assume what they contain without reading them first.

        The paths to the relevant files for this job are:
        - Job Script: {job_script_path}
        - Standard Output: {stdout_path}
        - Standard Error: {stderr_path}

        Follow these steps:
        1. **Start by reading the Standard Error file ({stderr_path})** as it is the most likely source of critical error messages.
        2. If the Standard Error file does not provide a clear explanation, **then read the Standard Output file ({stdout_path})**.
        3. **Also read the Job Script ({job_script_path})** to understand the intended commands.
        4. Synthesize your findings.

        Your final response MUST be structured with the following headings:

        **Problem**
        A one-sentence, high-level summary of the root cause of the job failure.

        **Explanation**
        A brief, clear, and non-technical explanation for the user, describing why the job failed.
        Reference specific errors or commands from the log files or script to support your explanation.
        Explain what the error means in simple terms.

        **Suggested Solution**
        Provide a specific, actionable step or command the user can take to fix the problem.
        If it's a script change, suggest the exact line or type of modification, along with the fixed script content using ({job_script_path}).
        If it's an environment issue, provide the command to set up the environment correctly (e.g., `module load`).
        Enclose any commands or code snippets in markdown code blocks.
        """),
        MessagesPlaceholder(variable_name="agent_scratchpad"),
    ])

    # Create the agent
    agent = create_openai_tools_agent(llm, tools, prompt)
    agent_executor = AgentExecutor(agent=agent, tools=tools, verbose=True)

    print(f"\n--- Starting analysis for Job ID: {job_id} ---")
    # Execute the agent with the job information
    result = agent_executor.invoke({
        "job_id": job_id,
        "job_script_path": job_script_path,
        "stdout_path": stdout_path,
        "stderr_path": stderr_path
    })

    # The LLM's final answer will be in result["output"]
    return result["output"]

# --- Main Execution Examples ---
if __name__ == "__main__":
    # add code to accept command line arguments for dir and job_id
    parser = argparse.ArgumentParser(description="Diagnose failed PBS Pro jobs.")
    parser.add_argument("job_id", help="Job ID to analyze")
    parser.add_argument("job_script_path", help="Path to the job script")
    parser.add_argument("stdout_path", help="Path to the standard output file")
    parser.add_argument("stderr_path", help="Path to the standard error file")
    args = parser.parse_args()
    if not args.job_script_path or not args.stdout_path or not args.stderr_path:
        parser.error("All file path arguments are required.")

    job_id = args.job_id
    job_script_path = args.job_script_path
    stdout_path = args.stdout_path
    stderr_path = args.stderr_path

    diagnosis = analyze_job_failure(job_id, job_script_path, stdout_path, stderr_path)

    print("\n" + "="*80 + "\n")
    print(f"--- Final Diagnosis for Job ID: {job_id} ---\n{diagnosis}")
    print("\n" + "="*80 + "\n")




