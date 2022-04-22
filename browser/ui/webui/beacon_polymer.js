
// This file is appended to polymer bundle
function beaconHideElement(content, selector) {
    const ele = content.querySelector(selector);
    if (!ele) {
        console.error(new Error('element not found ' + selector).stack); 
        return;
    }
    ele.style.display = 'none';
}

function applyBeaconGlobalStyle(content) {
    const styleElement = content.querySelector('style');
    if (!styleElement) {
        console.error(new Error('style element not found').stack); 
        return;
    }
    const attr = styleElement.getAttribute('include');
    if (!attr) attr = '';

    styleElement.setAttribute('include', attr + ' ' + 'beacon-global-style');
}

function registerBeaconGlobalStyle() {
    // Polymer style modules
    // https://polymer-library.polymer-project.org/3.0/docs/devguide/style-shadow-dom#style-modules
    const styleElement = document.createElement('dom-module');
    styleElement.innerHTML = 
 `<template>
   <style>
  
    /* Shared styles go here apply to elements as needed */
    .cr-nav-menu-item {
        border-end-end-radius: 4px;
        border-start-end-radius: 4px;
 
    }

    #urlCollectionToggle {
        display: none !important;
    } 
    
    html {
        --cr-card-border-radius: 4px;
        --cr-toolbar-search-field-border-radius: 4px;
    }
   </style>
 </template>`;

  styleElement.register('beacon-global-style');
}

const beaconOverrides = {
    "settings-personalization-options": function(content) {
        applyBeaconGlobalStyle(content);
    },
    "passwords-section": function(content) {
        beaconHideElement(content, '#manageLink');
    },
    "settings-payments-section": function(content) {
        beaconHideElement(content, '#manageLink');
    },
    "settings-security-page": function(content) {
        beaconHideElement(content, '#safeBrowsingEnhanced');
        beaconHideElement(content, '#safeBrowsingReportingToggle');
        beaconHideElement(content, '#passwordsLeakToggle');
        beaconHideElement(content, '#advanced-protection-program-link');
        beaconHideElement(content, '#banner');
    },
    "settings-privacy-page": function(content) {
        beaconHideElement(content, '#privacySandboxLinkRow');
    },
    "settings-site-settings-page": function(content) {
        beaconHideElement(content, '#banner');
    },
    "settings-menu": function(content) {
        applyBeaconGlobalStyle(content);
    }
 };
 // Based on trick used by Brave to override polymer elements
 const polymerElement__prepareTemplate = PolymerElement._prepareTemplate;
 PolymerElement._prepareTemplate = function() {
   polymerElement__prepareTemplate.call(this)
   const name = this.is
   if (!name ||
       !beaconOverrides.hasOwnProperty(name) ||
       !this.prototype || 
       !this.prototype._template ||
       !this.prototype._template.content) {
     return;
   }
 
   beaconOverrides[name](this.prototype._template.content);
 }

 registerBeaconGlobalStyle();
 