==============================
Introduction
==============================

Badaas enables the effortless construction of **distributed, resilient, 
highly available and secure applications by design**, while ensuring very simple 
deployment and management (NoOps).

.. warning::
   BaDaaS is still under development and each of its components can have a different state of evolution

Features and components
=================================

Badaas provides several key features, 
each provided by a component that can be used independently and has a different state of evolution:

- **Authentication** (unstable): Badaas can authenticate users using its internal 
  authentication scheme or externally by using protocols such as OIDC, SAML, Oauth2...
- **Authorization** (wip_unstable): On resource access, Badaas will check if the user 
  is authorized using a RBAC model.
- **Distribution** (todo): Badaas is built to run in clusters by default. 
  Communications between nodes are TLS encrypted using `shoset <https://github.com/ditrit/shoset>`_.
- **Persistence** (wip_unstable): Applicative objects are persisted as well as user files. 
  Those resources are shared across the clusters to increase resiliency. 
  To achieve this, BaDaaS uses the :doc:`badaas-orm <badaas-orm/index>` component.
- **Querying Resources** (unstable): Resources are accessible via a REST API.
- **Posix compliant** (stable): Badaas strives towards being a good unix citizen and 
  respecting commonly accepted norms. (see :doc:`badaas/configuration`)
- **Advanced logs management** (todo): Badaas provides an interface to interact with 
  the logs produced by the clusters. Logs are formatted in json by default.

Learn how to use BaDaaS following the :doc:`badaas/quickstart`.

.. toctree::
   :caption: BaDaaS

   self
   badaas/quickstart
   badaas/functionalities
   badaas/configuration

.. toctree::
   :caption: Badaas-orm

   badaas-orm/index
   badaas-orm/concepts
   badaas-orm/declaring_models
   badaas-orm/connecting_to_a_database
   badaas-orm/crud
   badaas-orm/query
   badaas-orm/advanced_query
   badaas-orm/preloading

.. toctree::
   :caption: Contributing

   contributing/contributing
   contributing/developing
   contributing/maintaining
   Github <https://www.github.com/ditrit/badaas>
