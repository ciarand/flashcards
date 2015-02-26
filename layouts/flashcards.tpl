{{define "flashcards"}}
<!doctype html>
<html>
<head>{{template "head" .}}</head>
<body>
    <div class="flowtime">
        {{range .Cards}}
        <div class="ft-section">
            <div class="ft-page">
                <div class="stack-center">
                    <p class="stacked-center">Question: {{.Question}}</p>
                </div>
            </div>
            <div class="ft-page">
                <div class="stack-center">
                    <p class="stacked-center">Answer: {{.Answer}}</p>
                </div>
            </div>
        </div>
        {{end}}
    </div>

    {{template "tail" .}}
</body>
</html>
{{end}}
