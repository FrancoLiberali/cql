==============================
Declaring models
==============================

Model declaration
-----------------------

The badaas-orm :ref:`model <badaas-orm/concepts:model>` declaration is based on the GORM model declaration, 
so its definition, conventions, tags and associations are compatible with badaas-orm. 
For details see `gorm documentation <https://gorm.io/docs/models.html>`_. 
On the contrary, badaas-orm presents some differences/extras that are explained in this section.

Base models
-----------------------

To be considered a model, your structures must have embedded one of the 
:ref:`base models <badaas-orm/concepts:base model>` provided by badaas-orm:

- `model.UUIDModel`: Model identified by a model.UUID (Random (Version 4) UUID).
- `model.UIntModel`: Model identified by a model.UIntID (auto-incremental uint).

Both base models provide date created, updated and `deleted <https://gorm.io/docs/delete.html#Soft-Delete>`_.

To use them, simply embed the desired model in any of your structs:

.. code-block:: go

  type MyModel struct {
    model.UUIDModel

    Name         string
    Email        *string
    Age          uint8
    Birthday     *time.Time
    MemberNumber sql.NullString
    ActivatedAt  sql.NullTime
    // ...
  }

Type of attributes
-----------------------

As we can see in the example in the previous section, 
the attributes of your models can be of multiple types, 
such as basic go types, pointers, and :ref:`nullable types <badaas-orm/concepts:nullable types>`.

This difference can generate differences in the information that is stored in the database, 
since saving a model created as follows:

.. code-block:: go

  MyModel{}

will save a empty string for Name but a null for the Email and the MemberNumber.

The use of nullable types is strongly recommended and badaas-orm takes into account 
their use in each of its functionalities.

Associations
-----------------------

All associations provided by GORM are supported.
For more information see <https://gorm.io/docs/belongs_to.html>, 
<https://gorm.io/docs/has_one.html>, <https://gorm.io/docs/has_many.html> and 
<https://gorm.io/docs/many_to_many.html>. 
However, in this section we will give some differences in badaas-orm and 
details that are not clear in this documentation.

IDs
^^^^^^^^^^^^^^^^^^^^^

Since badaas-orm base models use model.UUID or model.UIntID to identify the models, 
the type of id used in a reference to another model is the corresponding one of these two, 
for example:

.. code-block:: go

  type ModelWithUUID struct {
    model.UUIDModel
  }

  type ModelWithUIntID struct {
    model.UIntModel
  }

  type ModelWithReferences struct {
    model.UUIDModel

    ModelWithUUID *ModelWithUUID
    ModelWithUUIDID *model.UUID

    ModelWithUIntID *ModelWithUIntID
    ModelWithUIntIDID *model.UIntID
  }

References
^^^^^^^^^^^^^^^^^^^^^

References to other models can be made with or without pointers:

.. code-block:: go

  type ReferencedModel struct {
    model.UUIDModel
  }

  type ModelWithPointer struct {
    model.UUIDModel

    // reference with pointer
    PointerReference *ReferencedModel
    PointerReferenceID *model.UUID
  }

  type ModelWithoutPointer struct {
    model.UUIDModel

    // reference without pointer
    Reference ReferencedModel
    ReferenceID model.UUID
  }

As in the case of attributes, 
this can make a difference when persisting, since one created as follows:

.. code-block:: go

  ModelWithoutPointer{}

will also create and save an empty ReferencedModel{}, what may be undesired behavior. 
For this reason, although both options are still compatible with badaas-orm, 
we recommend the use of pointers for references. 
In case the relation is not nullable, use the `not null` tag in the id of the reference, for example:

.. code-block:: go

  type ReferencedModel struct {
    model.UUIDModel
  }

  type ModelWithPointer struct {
    model.UUIDModel

    // reference with pointer not null
    PointerReference *ReferencedModel
    PointerReferenceID *model.UUID `gorm:"not null"`
  }

Reverse reference
------------------------------------

Although no example within the `gorm's documentation <https://gorm.io/docs/has_one.html>`_ shows it, 
when defining relations, we can also put a reference in the reverse direction 
to add navigability to our model. 
In addition, adding this reverse reference will allow the corresponding conditions 
to be generated during condition generation.

For example:

.. code-block:: go

  type Related struct {
    model.UUIDModel

    YourModel *YourModel
  }

  type YourModel struct {
    model.UUIDModel

    Related *Related
    RelatedID *model.UUID
  }