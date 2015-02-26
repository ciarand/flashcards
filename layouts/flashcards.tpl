{{define "flashcards"}}
<!doctype html>
<html>
<head>{{template "head" .}}</head>
<body>
    {{template "before_content" .}}
    <p class="lead">
        It looks like you have {{len .Cards}} cards.

        <table class="table">
            <thead>
                <tr>
                    <th>Question</th>
                    <th>Answer</th>
                    <th>Class</th>
                </tr>
            </thead>
            <tbody>
                {{range .Cards}}
                <tr>
                    <td>{{.Question}}</td>
                    <td>{{.Answer}}</td>
                    <td>{{.Class}}</td>
                </tr>
                {{end}}
            </tbody>
        </table>
    </p>
    {{template "after_content" .}}
</body>
</html>
{{end}}
