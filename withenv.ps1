$ori = @{}
Try {
  $i = 0

  # Loading .env files
  if(Test-Path $args[0]) {
    foreach($line in (Get-Content $args[0])) {
      if($line -Match '^\s*$' -Or $line -Match '^#') {
        continue
      }

      $key, $val = $line.Split("=")
      $ori[$key] = if(Test-Path Env:\$key) { (Get-Item Env:\$key).Value } else { "" }
      New-Item -Name $key -Value $val -ItemType Variable -Path Env: -Force > $null
    }

    $i++
  }

  while(1) {
    if($i -ge $args.length) {
      exit
    }

    if(!($args[$i] -Match '^[^ ]+=[^ ]+$')) {
      break
    }

    $key, $val = $args[$i].Split("=")
    $ori[$key] = if(Test-Path Env:\$key) { (Get-Item Env:\$key).Value } else { "" }
    New-Item -Name $key -Value $val -ItemType Variable -Path Env: -Force > $null

    $i++
  }


  Invoke-Expression ($args[$i..$args.length] -Join " ")
} Finally {
  foreach($key in $ori.Keys) {
    New-Item -Name $key -Value $ori.Item($key) -ItemType Variable -Path Env: -Force > $null
  }
}