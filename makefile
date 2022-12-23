SHELL := C:\Users\Admin\AppData\Roaming\Microsoft\Windows\Start Menu\Programs\Windows PowerShell

run:
	$(shell cat .env | tr "\n" " ") go run .