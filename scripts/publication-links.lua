-- Local repository files are not bundled into PDF/EPUB. Keep chapter anchors
-- and external URLs intact, but send source/evidence links to the public repo.
function Link(link)
  local target = link.target
  if target:match('^#') or target:match('^[%a][%w+.-]*:') or target:match('^//') then
    return link
  end
  local path, suffix = target:match('^([^?#]*)(.*)$')
  local parts = target:match('^/') and {} or {'book'}
  for part in path:gmatch('[^/]+') do
    if part == '..' then
      table.remove(parts)
    elseif part ~= '.' then
      table.insert(parts, part)
    end
  end
  link.target = 'https://github.com/waywardgeek/ensemble/blob/main/' ..
    table.concat(parts, '/'):gsub(' ', '%%20') .. suffix
  return link
end
