import subprocess
import json
import sys

def main():
    # Prompt recibido desde Ollama
    prompt = " ".join(sys.argv[1:])
    
    # Aquí podemos hacer algo simple: convertir lenguaje natural a SQL
    # Este ejemplo es fijo para demo
    if "3 primeros usuarios" in prompt.lower():
        sql = "SELECT * FROM users LIMIT 3"
    else:
        sql = "SELECT * FROM users LIMIT 5"

    # Construimos payload MCP
    payload = {
        "version": "2.0",
        "type": "callTool",
        "tool": "query",
        "arguments": {"sql": sql},
        "callId": "1"
    }

    # Lanzamos el MCP server como subproceso
    proc = subprocess.Popen(
        ["go", "run", "mcp_mysql_query.go"],
        stdin=subprocess.PIPE,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE
    )

    # Enviamos el mensaje MCP serializado
    proc.stdin.write((json.dumps(payload) + "\n").encode())
    proc.stdin.flush()

    # Leemos la respuesta del MCP server
    response = proc.stdout.readline().decode().strip()
    print(response)

    # Cerramos el subproceso
    proc.terminate()

if __name__ == "__main__":
    main()

