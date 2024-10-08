# Project Directory Structure

```
alpha-golang/
├── .git/
    ├── hooks/
    │   ├── applypatch-msg.sample
    │   ├── commit-msg.sample
    │   ├── fsmonitor-watchman.sample
    │   ├── post-update.sample
    │   ├── pre-applypatch.sample
    │   ├── pre-commit.sample
    │   ├── pre-merge-commit.sample
    │   ├── pre-push.sample
    │   ├── pre-rebase.sample
    │   ├── pre-receive.sample
    │   ├── prepare-commit-msg.sample
    │   ├── push-to-checkout.sample
    │   ├── sendemail-validate.sample
    │   └── update.sample
    ├── info/
    │   └── exclude
    ├── logs/
    │   ├── refs/
    │   │   ├── heads/
    │   │   │   └── master
    │   │   └── remotes/
    │   │   │   └── origin/
    │   │   │       ├── HEAD
    │   │   │       └── master
    │   └── HEAD
    ├── objects/
    │   ├── 01/
    │   │   └── 4b3b3c69bc12e21a0e729fb817139e5c4d6c1f
    │   ├── 02/
    │   │   └── 3a5bd37caf778780a49911af685f3d23d28af9
    │   ├── 03/
    │   │   └── 4a8e642ee007f8b1080ffe12a3340e937e56ec
    │   ├── 05/
    │   │   └── bd63595bdef370b7783104a840b37ca6e1f280
    │   ├── 09/
    │   │   └── e3721a6b6702dd0e684ba5c1c8f6b5fea61fc7
    │   ├── 0a/
    │   │   └── db25da359003363c3f2d9d0570faf63563300e
    │   ├── 0b/
    │   │   └── 7d9f06673699f0727c7d6580d089525b6653c8
    │   ├── 0f/
    │   │   └── dc967309c6bbf7fc7f70ac0e64262e07d8c06e
    │   ├── 10/
    │   │   └── a10220ac517f00e27e43f466ef7d17c1a672ac
    │   ├── 12/
    │   │   ├── a284143ab07704d7b061dc6581c5d2b7b7b5ce
    │   │   └── aa9f92de6e64f85536b4681cff06eea4da89ad
    │   ├── 15/
    │   │   └── 8f26afdac380a791bc8710c2bfaabe82f50933
    │   ├── 16/
    │   │   ├── 27fb17dfa9f462f47fa150968f93f6f3a75283
    │   │   └── 9e0c527359d8e32517f455e8a5d8b8462a7901
    │   ├── 20/
    │   │   └── 9816c03cbe70d1a3b319f9cdca988214904ceb
    │   ├── 21/
    │   │   └── bbfe68adf03ca3e57e262cdf28980903ac6703
    │   ├── 22/
    │   │   └── 96d61fad0b34bda13f779905371e24b3c6dea6
    │   ├── 23/
    │   │   └── ec989a3112f3bc5a5eaca75356f71a1a509a37
    │   ├── 25/
    │   │   └── 43829d6e74cc7e8a56123b5ecae513410947d4
    │   ├── 26/
    │   │   └── aea326d59ee9ec3dd335babfa7f47cd9da8f7f
    │   ├── 27/
    │   │   ├── 55e820b22413eb91508d7a82ab2f514fac8450
    │   │   └── 9be4641bd7f5215e115b87930f8c59e22c2f4e
    │   ├── 29/
    │   │   └── 96eb7c866f96f7ee83b8233cf3fe7e3bfa0d40
    │   ├── 2a/
    │   │   └── 87bfbc0d99f1035385729deedd8bed6bb62be0
    │   ├── 30/
    │   │   ├── 5af05a234f1d8c0c9b1e48d038589c787ef431
    │   │   ├── df2729325830dc87f5c72649868e34381de6ca
    │   │   └── f336201a45ee0dce2a0e3ece4594a848e8e654
    │   ├── 32/
    │   │   └── 16e099f54ca3badc6af489b430f3b402ce4993
    │   ├── 34/
    │   │   └── ae576d9482efff277ac80a02b6d522bd24ad8b
    │   ├── 36/
    │   │   ├── 9bf5e51bf5447eaa6ee7293c2b3a896670573e
    │   │   └── ee1db529e2d6c8257d8a3cab28a6fd5af28ae5
    │   ├── 3a/
    │   │   └── 97243c21ef34537289cb834c239c20c5fc4d1f
    │   ├── 44/
    │   │   └── abc0c8560d5f4dfbdc489d3c86b30aba082de2
    │   ├── 4b/
    │   │   └── e69d83bc1f5a195804e02746b3aead49a06f90
    │   ├── 4c/
    │   │   └── 0ee0b1d78ff737444d842dc4334ff6ceacb4ca
    │   ├── 50/
    │   │   └── c69819e8ec8ef72ce1656e0602062f7bfcf0fb
    │   ├── 54/
    │   │   └── ff196eead67f1b6eae78325ee644dcc65d4400
    │   ├── 55/
    │   │   └── 74f1fcd793027ff9851d85e7d702009c571509
    │   ├── 56/
    │   │   └── 6ad0cc4c56388c6f6b4867d7438ae1f013b77d
    │   ├── 57/
    │   │   └── 10ca1fcf972647a9e16d0b4deb54f30d37514c
    │   ├── 58/
    │   │   └── 92cb63106567be76caa14286a707995fd0add6
    │   ├── 5c/
    │   │   └── 01726e72c74f2dbc917a9cf8f0d5f4b5fe18d6
    │   ├── 5e/
    │   │   └── 764c4f0b9a64bb78a5babfdd583713b2df47bf
    │   ├── 5f/
    │   │   ├── 919b1149983424e9eaa3a5bc976261e67d2419
    │   │   └── 94bd5d6041bd34296565c97a8d900ddf0d9c1e
    │   ├── 66/
    │   │   └── a52f3c802b439bb1047ff0be71ab629c843c9a
    │   ├── 6a/
    │   │   └── c7d22a51aa6b03ad3c25861401714ec4284490
    │   ├── 6b/
    │   │   └── 7a8e972a63f057d80fa8d36fe9b9ea8c4f418b
    │   ├── 71/
    │   │   └── 9626eb4911f732900fbe0b8fdef71240d04b6a
    │   ├── 73/
    │   │   └── 981b0b7d6910f0cd22085b8da73a8f9621090f
    │   ├── 74/
    │   │   └── 15643b747526215884deb28333abb760a1edd4
    │   ├── 77/
    │   │   └── 31c5e6fc7d94b476d62d556b6a52e154da61ed
    │   ├── 78/
    │   │   └── 31a3ebf299445ebc4e1838247ea2347dfe57b7
    │   ├── 82/
    │   │   └── 8446f8ee1991eee6a15ef8042c8cb66496906d
    │   ├── 83/
    │   │   └── 0a5346ffa81362e7fe1720240d22ac477ff28a
    │   ├── 86/
    │   │   └── efe6ab1b0f6cf6d75cc9ece8fada1f2d6a2e05
    │   ├── 8b/
    │   │   └── 56db0f4f0f9b69e97fb3b038c026e987b07c65
    │   ├── 92/
    │   │   └── 8100fec49547443d9d6c7a3cdbe64a13c8991b
    │   ├── 93/
    │   │   └── 36417885507721033ee75d32584ee8eff0da3d
    │   ├── 94/
    │   │   └── 6ff6e77b3638b0d9de61306d09f34439cba29e
    │   ├── 95/
    │   │   └── fee72ed9d37ebf11be750ac69a5f7620ee09b0
    │   ├── 98/
    │   │   └── 45ad5e7307f9ca7726ea9505e29066e060d42e
    │   ├── 9a/
    │   │   └── 1086d6ab96ff2e3e774f209392fda5b58e5f0b
    │   ├── 9d/
    │   │   └── fc7879df42bae2f0cfb9af130a2bc1b2947d75
    │   ├── a2/
    │   │   ├── 5e7c897c3fe338ebfbf8d1c9118c910057868e
    │   │   └── 781ec881ee993052489bd228893cb5e3b8e7cd
    │   ├── a9/
    │   │   └── f20de1f0ca3c0b0eef70aea81bc1d19896e1a2
    │   ├── ac/
    │   │   └── ba91db8fc2797554550d7256e7a63788c99ef2
    │   ├── ae/
    │   │   └── b656a2389acc7109e5bd807ed23b52b5700fa9
    │   ├── af/
    │   │   └── 668fb192a9e33429fc287ecb0e639b90794c68
    │   ├── b1/
    │   │   └── 616f4191d1298c6d5c70872b74644fd8101ab6
    │   ├── b3/
    │   │   └── c74f036e67a20e4b7834c854afe8b2b7aa6fad
    │   ├── b8/
    │   │   ├── 8c1341d8f430d54be16050d901b78ed8a39aae
    │   │   └── af426606d2daa11ccf7e8983ffa3df3c493f87
    │   ├── b9/
    │   │   └── 9fbb9b169bb68cc2221fe0b71cc6ed02f1c136
    │   ├── ba/
    │   │   └── c9ea409c5a3b3a6ebbb2be44f7e9b677be3736
    │   ├── bc/
    │   │   └── 771b8d790e29d3d6c4442c335fd2bf83a6765c
    │   ├── bd/
    │   │   └── 5e31442ca6cf544984d5aaa76202df1226f25d
    │   ├── c2/
    │   │   └── 55c9db3d23496e1ec4d7c38f1ab8345f484030
    │   ├── c4/
    │   │   └── 15bb7c6be2f2c8bde821bbe07cd006042abfca
    │   ├── c7/
    │   │   └── 953b4a158127b900af0b613bf5179c75c8ae5d
    │   ├── cb/
    │   │   └── a216a2186fa7b7c7402d6c2fbdc9f0320c3be9
    │   ├── ce/
    │   │   └── da188a01de9aa857a8a4c09d2c397f46e9b5c6
    │   ├── d4/
    │   │   ├── 073c2fdc89587aa4010bf82d6a8239aecdc863
    │   │   └── a432af30a2d55bdc85b6f94e174e47a16e677c
    │   ├── d5/
    │   │   └── a1d092d6eb9b32cf0f20ef14070888daad59a1
    │   ├── d7/
    │   │   └── 803d0aa3d7168d5dee951f3dea95f8bdee8b0f
    │   ├── d9/
    │   │   ├── 07dca44bfe4f7081e66f326e2f63d125372416
    │   │   └── 75886635adf8c7d0391d2d382907ad55f8af99
    │   ├── de/
    │   │   └── 633a2ea94cbbaad78e29a04dc24bf501d0a24e
    │   ├── df/
    │   │   ├── 573701cb6af6d1064417a5d9f8d3aaa123e3dc
    │   │   └── dd76761bb785f5cc20310727ddbd80917fd1e5
    │   ├── e0/
    │   │   ├── dfed3157778e2bb27c97a2b166ae86a66ca099
    │   │   └── e25abb1e76df4a243c9e9ac3f871aeac3a8583
    │   ├── e1/
    │   │   └── f4ecce4c09b8cbe6a5b2839adf7001c6e1c330
    │   ├── e6/
    │   │   ├── 1ff54a1538cf6b40244aa09c5abd7ca6f92065
    │   │   └── 76d216f5a517bc02ab13170dbe56cb991ba8bc
    │   ├── e9/
    │   │   └── 730561656ee74b03731a118478c975c566b85a
    │   ├── eb/
    │   │   └── 54876c0acef60a39b94cebc2a76a6801c57e36
    │   ├── f1/
    │   │   └── 450043fedf4229b7e20692879140785f35a0ec
    │   ├── f2/
    │   │   └── 842ebfa1229ea4fd52d60a1d3cc2dd55ea670b
    │   ├── fa/
    │   │   └── f94cdde860999b8c5a37b5ef79562756918dc7
    │   ├── fb/
    │   │   ├── df9d78f5067148333d30e951148da374bd1163
    │   │   └── f0422da39dcc9cb68936d3b6db5e28efb071d2
    │   ├── fd/
    │   │   └── 4648eeeb2adfac70199dde235a405c1c8f0b08
    │   ├── ff/
    │   │   └── 04bbd483c902acbc942d40be8b4c399c98e60e
    │   ├── info/
    │   └── pack/
    ├── refs/
    │   ├── heads/
    │   │   └── master
    │   ├── remotes/
    │   │   └── origin/
    │   │   │   ├── HEAD
    │   │   │   └── master
    │   └── tags/
    │   │   ├── v0.0.1
    │   │   ├── v0.0.2
    │   │   └── v0.0.3
    ├── COMMIT_EDITMSG
    ├── FETCH_HEAD
    ├── HEAD
    ├── ORIG_HEAD
    ├── config
    ├── description
    └── index
├── .idea/
    ├── .gitignore
    ├── alpha-golang.iml
    ├── modules.xml
    ├── vcs.xml
    └── workspace.xml
├── api/
    ├── .env
    ├── database.go
    ├── go.mod
    ├── go.sum
    ├── handlers.go
    ├── main.go
    └── middleware.go
├── bot/
    ├── bin/
    │   └── discordbot
    ├── .env
    ├── bot.exe
    ├── database.go
    ├── digests.go
    ├── go.mod
    ├── go.sum
    ├── handlers.go
    ├── llm.go
    ├── main.go
    ├── openai.go
    ├── openai_test.go
    ├── printProject.py
    ├── summarize.go
    ├── synthesize.go
    └── talk.go
├── common/
    ├── go.mod
    └── models.go
├── database/
    └── migrations/
    │   ├── 001_create_messages_table.up.sql
    │   ├── 002_create_summaries_table.up.sql
    │   ├── 003_create_digests_table.up.sql
    │   ├── 004_add_summarized_column_to_messages.up.sql
    │   ├── 005_add_digested.up.sql
    │   ├── 006_add_synthesized.up.sql
    │   ├── 007_add_synthesizedTEXT.up.sql
    │   └── README.md
├── .gitignore
├── printProject.py
├── project_structure.md
└── railway.toml
```

# File Contents

