// Modified for Beacon use
/* Copyright (c) 2021 The Beacon Authors. All rights reserved.
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this file,
 * You can obtain one at http://mozilla.org/MPL/2.0/. */

#include "base/containers/contains.h"
#include "base/strings/string_split.h"
#include "third_party/re2/src/re2/re2.h"

#include "src/chrome/browser/ui/webui/about_ui.cc"

std::string AboutUIHTMLSource::ChromeURLs() const {
  std::string chrome_urls = ::ChromeURLs();

  // Replace Chrome -> Beacon.
  const std::string chrome_header = "Chrome URLs";
  const std::string beacon_header = "Beacon URLs";
  const std::string chrome_pages_header = "List of Chrome URLs";
  const std::string beacon_pages_header = "List of Beacon URLs";
  const std::string chrome_internal_pages_header =
      "List of chrome://internals pages";
  const std::string beacon_internal_pages_header =
      "List of beacon://internals pages";
  const std::string chrome_url_list = ">chrome://";
  const std::string beacon_url_list = ">beacon://";
  RE2::GlobalReplace(&chrome_urls, chrome_header, beacon_header);
  RE2::GlobalReplace(&chrome_urls, chrome_pages_header, beacon_pages_header);
  RE2::GlobalReplace(&chrome_urls, chrome_internal_pages_header,
                     beacon_internal_pages_header);
  RE2::GlobalReplace(&chrome_urls, chrome_url_list, beacon_url_list);

  // Remove some URLs.
  auto html_lines = base::SplitStringPiece(
      chrome_urls, "\n", base::KEEP_WHITESPACE, base::SPLIT_WANT_ALL);
  const base::flat_set<base::StringPiece> kURLsToRemove{
      "beacon://memories",
  };
  // URLs in html should be sorted so it's okay to iterate over sorted
  // kURLsToRemove.
  auto html_line_it = html_lines.begin();
  auto url_to_remove_it = kURLsToRemove.begin();
  while (html_line_it != html_lines.end() &&
         url_to_remove_it != kURLsToRemove.end()) {
    if (base::Contains(*html_line_it, *url_to_remove_it)) {
      html_line_it = html_lines.erase(html_line_it);
      ++url_to_remove_it;
    } else {
      ++html_line_it;
    }
  }

  return base::JoinString(html_lines, "\n");
}
