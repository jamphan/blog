Remove-Item .\public\* -Recurse -Force
.\hugoextended.ps1
gsutil rsync -R .\public\ gs://www.jamphan.dev