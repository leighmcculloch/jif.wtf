configure :build do
  activate :asset_hash
end

activate :s3_sync do |s3_sync|
  s3_sync.bucket                     = 'jif.wtf'
  s3_sync.region                     = 'us-east-1'
  s3_sync.delete                     = false
  s3_sync.after_build                = false
  s3_sync.prefer_gzip                = true
  s3_sync.path_style                 = true
end
