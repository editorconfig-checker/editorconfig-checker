@for /r . %%A in (*.txt *.html) do @(
    printf "%%-100s: " "%%~A"
    uchardet "%%~A"
)
