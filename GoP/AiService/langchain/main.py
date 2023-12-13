import requests

res = requests.post('http://127.0.0.1:11434/api/chat', data={"model":"orca-mini","messages":[{"role":"user","content":"Write 12 word start on leter S"}]})
print(res.text)
