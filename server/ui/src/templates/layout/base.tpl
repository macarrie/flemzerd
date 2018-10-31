<!doctype html>
<html lang="en">
    <head>
        <meta charset="utf-8">
        <title>Flemzerd</title>
        <base href="/">

        <meta name="viewport" content="width=device-width, initial-scale=1">
        <meta name="turbolinks-cache-control" content="no-cache">

        <link rel="icon" type="image/x-icon" href="favicon.ico">
        <link href="https://fonts.googleapis.com/icon?family=Material+Icons" rel="stylesheet">
        <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/font-awesome/4.7.0/css/font-awesome.min.css">
        <link rel="stylesheet" href="static/css/style.css">

        <script src="static/js/bundle.js" data-turbolinks-suppress-warning></script>
        <script src="https://cdnjs.cloudflare.com/ajax/libs/uikit/3.0.0-rc.20/js/uikit.min.js"></script>
        <script src="https://cdnjs.cloudflare.com/ajax/libs/uikit/3.0.0-rc.20/js/uikit-icons.min.js"></script>
    </head>
    <body>

        {{ include "layout/header" }}

        {{ template "content" . }}

        {{ include "layout/footer" }}
    </body>
</html>
