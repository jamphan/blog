Remove-Item .\public\* -Recurse -Force
.\hugoextended.ps1 -D
gsutil rsync -R .\public\ gs://www.jamphan.dev