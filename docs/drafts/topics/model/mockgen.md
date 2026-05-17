

```bash
mockgen \
  -destination ./examplemodel_mock.go \
  -package model \
  . \
  ExampleModel
```  



```bash
mockgen \
  -destination examplemodel_mock.go \
  -package model \
  examplemodel_gen.go \
  ExampleModel
```

